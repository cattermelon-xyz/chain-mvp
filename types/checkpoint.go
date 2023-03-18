package types

import (
	"log"
	"math"

	"github.com/hectagon-finance/chain-mvp/third_party/tree"
	"github.com/hectagon-finance/chain-mvp/third_party/utils"
	"github.com/hectagon-finance/chain-mvp/types/event"
)

const NoFallbackOption = math.MaxUint64 - 1
const EndOfMission = math.MaxUint64 - 2
const GenesisBlock = 0

type CheckPointStartedStatus string

const (
	ChkPFailToStart CheckPointStartedStatus = "CheckPoint Fail To Start"
	ChkPStarted     CheckPointStartedStatus = "CheckPoint Started"
	ChkPIsAnOutput  CheckPointStartedStatus = "Output Event Emitted"
	ChkPNil         CheckPointStartedStatus = "CheckPoint is nil"
)

type CheckPoint struct {
	Id               string
	Title            string
	Description      string
	FallbackId       uint64
	children         []*CheckPoint
	voteMachine      VotingMachine
	lastBlockToVote  uint64
	lastBlockToTally uint64
	outputEvent      *event.Event
	mission          *Mission
}

// return something that is printable
func (n *CheckPoint) Data() interface{} {
	return n.Title
}

// cannot return n.children directly.
// https://github.com/golang/go/wiki/InterfaceSlice
func (n *CheckPoint) Children() (c []tree.Node) {
	for _, child := range n.children {
		c = append(c, tree.Node(child))
	}
	return
}

/**
* Return an Output with Eventa.
* When this node start(), an event will be emitted.
 */
func (mission *Mission) CreateOutput(name string, desc string, ev *event.Event) *CheckPoint {
	c := CheckPoint{
		Id:               utils.RandString(16),
		Title:            name,
		Description:      desc,
		FallbackId:       NoFallbackOption,
		children:         nil,
		voteMachine:      nil,
		lastBlockToVote:  GenesisBlock,
		lastBlockToTally: GenesisBlock,
		outputEvent:      ev,
		mission:          mission,
	}
	return &c
}
func (this *CheckPoint) Attach(childId string) *CheckPoint {
	child := GetCheckPointById(childId)
	if this.children == nil {
		this.children = make([]*CheckPoint, 0)
	}
	if this.mission != child.mission {
		log.Println("Mismatch mission id")
		return nil
	}
	this.children = append(this.children, child)
	return child
}

/**
* Conversational text the describe the current state of the CheckPoint
* including: Title, Description, Options, How voting will conduct
**/
func (this *CheckPoint) Print() {
	log.Printf("%s\n%s\nVoting Mechanism:\n%s\n", this.Title, this.Description, this.voteMachine.Desc())
	for i := range this.children {
		log.Printf("- opt %d: %s\n", i, this.children[i].Title)
	}
	log.Printf("\n")
}
func (this *CheckPoint) Get(idx uint64) *CheckPoint {
	if idx < uint64(len(this.children)) {
		return this.children[idx]
	}
	return nil
}

func (this *CheckPoint) GetMissionId() Address {
	return this.mission.id
}

func (this *CheckPoint) GetOutputEvent() *event.Event {
	return this.outputEvent
}

/**
* Return CheckPointStartedStatus: {OutputEventEmitted, ChkPFailToStart, ChkPStarted}
 */
func (this *CheckPoint) start(lastTalliedResult []byte) (CheckPointStartedStatus, *event.Event) {
	// an Output or a CheckPoint with votingMachine?
	log.Printf("%s\n", this.Title)
	if this.outputEvent != nil {
		return ChkPIsAnOutput, this.outputEvent
	} else { // a CheckPoint with votingMachine
		if this.children == nil {
			log.Printf("func (this *CheckPoint) start([]byte): this.children is nil\n")
			return ChkPFailToStart, nil
		}
		if len(this.children) == 0 || this.FallbackId == NoFallbackOption {
			log.Printf("func (this *CheckPoint) start([]byte): len(this.children) == 0 or this.FallbackId == NoFallbackOption\n")
			return ChkPFailToStart, nil
		}
	}
	if this.voteMachine == nil {
		log.Printf("func (this *CheckPoint) start([]byte): this.voteMachine is nil\n")
		return ChkPFailToStart, nil
	}
	started := this.voteMachine.Start(lastTalliedResult, uint64(len(this.children)), this.FallbackId)
	if started == true {
		return ChkPStarted, nil
	}
	return ChkPFailToStart, nil
}
func (this *CheckPoint) isValidChoice(option []byte) bool {
	if this.voteMachine.IsStarted() == false {
		return false
	}
	return this.voteMachine.ValidateVote(option)
}

/*
Id               string
Title            string
Description      string
FallbackId       uint64
ChildrenId       []string
LastBlockToVote  uint64
LastBlockToTally uint64
OutputEventId    string
OutputEventName  string
OutputEventArgs  []byte
VoteMachineType  string
VoteMachine      []byte
*/
func (this *CheckPoint) marshal() CheckPointData {
	ChildrenId := make([]string, 0)
	if this.children != nil {
		for _, chkp := range this.children {
			ChildrenId = append(ChildrenId, chkp.Id)
		}
	}
	OutputEventId := ""
	OutputEventName := ""
	OutputEventArgs := make([]byte, 0)
	if this.outputEvent != nil {
		OutputEventId = this.outputEvent.Id
		OutputEventName = this.outputEvent.Name
		OutputEventArgs = this.outputEvent.Args
	}
	return CheckPointData{
		Id:               this.Id,
		Title:            this.Title,
		Description:      this.Description,
		FallbackId:       this.FallbackId,
		ChildrenId:       ChildrenId,
		LastBlockToVote:  this.lastBlockToVote,
		LastBlockToTally: this.lastBlockToTally,
		OutputEventId:    OutputEventId,
		OutputEventName:  OutputEventName,
		OutputEventArgs:  OutputEventArgs,
	}
}

/**
* Function vote
* Params: tr *Mission, who string, input []byte
* Returns: recordStatus ExecutionStatus, tallyStatus ExecutionStatus, newChkPStatus ExecutionStatus, fallbackAttempt bool
* TODO: what if we want to hide the voter's option from validator?
 */
func (this *CheckPoint) vote(who string, input []byte) (ExecutionStatus, ExecutionStatus, ExecutionStatus, bool, *event.Event) {
	var recordStatus ExecutionStatus = DIDNOTSTART
	var tallyStatus ExecutionStatus = DIDNOTSTART
	var newChkPStatus ExecutionStatus = DIDNOTSTART
	var ev *event.Event = nil
	fallbackAttempt := false
	// check for fallback
	fallbackAttempt, newChkPStatus, ev = fallback(this.mission, this.voteMachine, this.FallbackId)
	if fallbackAttempt == false {
		if this.voteMachine.Record(who, input) == true {
			recordStatus = SUCCEED
		} else {
			recordStatus = FAILED
		}
		// then check for tally
		if fallbackAttempt == false && this.voteMachine.ShouldTally() == true {
			tallyStatus, newChkPStatus, ev = tally(this.mission, this.voteMachine)
		}
	}
	return recordStatus, tallyStatus, newChkPStatus, fallbackAttempt, ev
}

/**
* Func tally; count all the vote
* Args: tr *Mission, m VotingMachine, input []byte
* Return: tallyStatus ExecutionStatus, newChkPoinStatus ExecutionStatus
 */
func tally(mission *Mission, m VotingMachine) (ExecutionStatus, ExecutionStatus, *event.Event) {
	_tallyStatus, _, tallyResult, selectedOption := m.Tally()
	_newChkPointStatus := ChkPFailToStart
	tallyStatus := FAILED
	newChkPointStatus := DIDNOTSTART
	var outputEvent *event.Event = nil
	if _tallyStatus == true {
		tallyStatus = SUCCEED
		if selectedOption != NoOptionMade {
			_newChkPointStatus, outputEvent = mission.Choose(selectedOption, tallyResult)
			// fmt.Println("tally: ", _newChkPointStatus)
			if _newChkPointStatus == ChkPStarted || _newChkPointStatus == ChkPIsAnOutput {
				newChkPointStatus = SUCCEED
			} else {
				newChkPointStatus = FAILED
			}
		}
	}
	return tallyStatus, newChkPointStatus, outputEvent
}

/**
* Func fallback; check if Voter can no longer vote, Mission can no longer tally then choose fallbackId
* Args: tr *Mission, m VotingMachine, fallbackId uint64, input []byte
* Return: fallbackAttempt bool, newChkPointStatus ExecutionStatus
 */
func fallback(mission *Mission, m VotingMachine, fallbackId uint64) (bool, ExecutionStatus, *event.Event) {
	currentBlk := mission.blockchain.GetCurrentBlockNumber()
	lastBlkVote := mission.currentChkP.lastBlockToVote
	lastBlkTally := mission.currentChkP.lastBlockToTally
	tallyResult, selectedOption := m.GetTallyResult()
	newChkPointStatus := DIDNOTSTART
	var ev *event.Event = nil
	if currentBlk > lastBlkVote && currentBlk > lastBlkTally && selectedOption == NoOptionMade {
		_newChkPointStatus, ev := mission.Choose(fallbackId, tallyResult)
		// fmt.Println("fallback: ", _newChkPointStatus)
		if _newChkPointStatus == ChkPStarted || _newChkPointStatus == ChkPIsAnOutput {
			newChkPointStatus = SUCCEED
		} else {
			newChkPointStatus = FAILED
		}
		return true, newChkPointStatus, ev
	}
	return false, newChkPointStatus, ev
}
