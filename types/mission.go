package types

import (
	"errors"

	"github.com/hectagon-finance/chain-mvp/third_party/tree"
	"github.com/hectagon-finance/chain-mvp/third_party/utils"
	"github.com/hectagon-finance/chain-mvp/types/event"
)

type ExecutionStatus uint8

const (
	DIDNOTSTART ExecutionStatus = iota
	FAILED
	SUCCEED
)

type Mission struct {
	id       Address
	Title    string
	Fulltext string
	Owner    Address

	blockchain  Blockchain
	startChkP   *CheckPoint
	currentChkP *CheckPoint
	isStarted   bool
	isActive    bool
}

var Missions = make([]*Mission, 0)
var cachedMissions = make(map[Address]*Mission)

func GetMissionById(id Address) (*Mission, error) {
	cached := cachedMissions[id] // cache hit
	if cached != nil {
		for _, m := range Missions {
			if m.id == id {
				cachedMissions[id] = m
				return m, nil
			}
		}
	}
	return nil, errors.New(string(id) + " not found")
}

var CheckPoints = make([]*CheckPoint, 0)
var cachedCheckPoints = make(map[string]*CheckPoint)

func GetCheckPointById(id string) *CheckPoint {
	cached := cachedCheckPoints[id] // cache hit
	if cached != nil {
		for _, chkp := range CheckPoints {
			if chkp.Id == id {
				cachedCheckPoints[id] = chkp
				return chkp
			}
		}
	}
	return nil
}

func CreateMission(title string, fulltext string, b Blockchain) (*Mission, string) {
	id := utils.RandString(16)
	i := Mission{
		id:          Address(id),
		Title:       title,
		Fulltext:    fulltext,
		startChkP:   nil,
		currentChkP: nil,
		isStarted:   false,
		isActive:    false,
		blockchain:  b,
	}
	Missions = append(Missions, &i)
	return &i, id
}

func (mission *Mission) CreateEmptyCheckPoint(title string, desc string, b VotingMachine) *CheckPoint {
	CheckPoint := CheckPoint{
		Id:               utils.RandString(16),
		Title:            title,
		Description:      desc,
		children:         []*CheckPoint{},
		voteMachine:      b,
		FallbackId:       NoFallbackOption,
		lastBlockToVote:  GenesisBlock,
		lastBlockToTally: GenesisBlock,
		outputEvent:      nil,
		mission:          mission,
	}
	return &CheckPoint
}
func (mission *Mission) CreateCheckPoinWithChildren(name string, desc string, children []*CheckPoint, b VotingMachine, fallbackId uint64, lastBlockToVote uint64, lastBlockToTally uint64) *CheckPoint {
	c := CheckPoint{
		Id:               utils.RandString(16),
		Title:            name,
		Description:      desc,
		FallbackId:       fallbackId,
		children:         children,
		voteMachine:      b,
		lastBlockToVote:  lastBlockToVote,
		lastBlockToTally: lastBlockToTally,
		outputEvent:      nil,
		mission:          mission,
	}
	return &c
}

func (this *Mission) GetId() Address {
	return this.id
}

// TODO: is it safe to do this? should we check all the nodes and events (observer)?
func DeleteMission(id string) bool {
	found := -1
	for idx, m := range Missions {
		if string(m.id) == id {
			found = idx
		}
	}
	if found != -1 {
		Missions[found] = Missions[len(Missions)-1]
		Missions = Missions[:len(Missions)-1]
		return true
	}
	return false
}
func (this *Mission) SetStartChkP(chkPId string) {
	chkP := GetCheckPointById(chkPId)
	this.startChkP = chkP
}
func (this *Mission) SetCurrentChkP(chkPId string) {
	chkP := GetCheckPointById(chkPId)
	this.currentChkP = chkP
}
func (this *Mission) CurrentChkP() *CheckPoint {
	return this.currentChkP
}
func (this *Mission) StartChkP() *CheckPoint {
	return this.startChkP
}

// func (this *Mission) edit(d Mission) bool {
// 	return false
// }

/**
* Function Start
 */
func (this *Mission) Start() bool {
	if this.isStarted == false && this.startChkP != nil {
		nodeStarted, outputEvent := this.startChkP.start(nil)
		if nodeStarted == ChkPStarted {
			this.isStarted = true
			this.isActive = true
			this.SetCurrentChkP(this.startChkP.Id)
			ev := this.blockchain.GetEventManager()
			ev.EmitMissionStarted(string(this.id))
			if outputEvent != nil {
				ev.Emit(outputEvent.Id)
			}
		} else {
			// log.Printf("func (this *Mission) Start(): Mission cannot start with nodeStarted is [%s]\n", nodeStarted)
		}
	}
	return this.isStarted
}

func (this *Mission) Stop() {
	if this.isStarted == true {
		this.isStarted = false
		this.isActive = false
		ev := this.blockchain.GetEventManager()
		ev.EmitMissionStopped(string(this.id))
	}
}

func (this *Mission) Pause() {
	if this.isActive == true && this.isStarted == true {
		this.isActive = false
		ev := this.blockchain.GetEventManager()
		ev.EmitMissionPaused(string(this.id))
	}
}

func (this *Mission) Resume() (bool, error) {
	if this.isStarted == true {
		this.isActive = true
		ev := this.blockchain.GetEventManager()
		ev.EmitMissionResumed(string(this.id))
		return true, nil
	}
	return false, errors.New("func (this *Mission) Resume()" + string(this.id) + " is stopped, can not start again")
}

func (this *Mission) PrintFromStart() {
	tree.Print(this.startChkP)
}
func (this *Mission) PrintFromCurrent() {
	if this.isStarted != true {
		tree.Print(this.startChkP)
	} else {
		tree.Print(this.currentChkP)
	}
}

func (this *Mission) Choose(idx uint64, tallyResult []byte) (CheckPointStartedStatus, *event.Event) {
	nextChkP := this.currentChkP.Get(idx)
	var started CheckPointStartedStatus
	var ev *event.Event = nil
	if nextChkP == nil {
		// log.Println("func (this *Mission) Choose(uint64, []byte): " + strconv.FormatUint(idx, 10) + " out of bound, no move")
		started = ChkPNil
		this.Stop()
	} else {
		started, ev = nextChkP.start(tallyResult)
		// log.Printf("from %s choose: %d got %s %s\n", this.currentChkP.Title, idx, nextChkP.Title, started)
		if started == ChkPStarted {
			this.currentChkP = nextChkP
		} else if started == ChkPIsAnOutput {
			this.Stop()
		}
	}
	return started, ev
}
func (this *Mission) IsValidChoice(option []byte) bool {
	return this.currentChkP.isValidChoice(option)
}

/**
* Function Vote
* Params: option []byte, who string, checkPointId string
* Returns: voteRecordedSucceed ExecutionStatus, talliedSucceed ExecutionStatus, newChkPointStartedSucceed ExecutionStatus, fallbackAttempt bool
 */
func (this *Mission) Vote(option []byte, who string, checkPointId string) (ExecutionStatus, ExecutionStatus, ExecutionStatus, bool, *event.Event) {
	ev := event.EventManager(this.blockchain.GetEventManager())
	voteRecordStatus := DIDNOTSTART
	tallyStatus := DIDNOTSTART
	newChkPointStatus := DIDNOTSTART
	fallbackAttempt := false
	var outputEvent *event.Event = nil
	if this.isActive == false {
		return DIDNOTSTART, DIDNOTSTART, DIDNOTSTART, false, outputEvent
	}
	lastChkPointId := this.currentChkP.Id
	if this.IsValidChoice(option) == true && this.currentChkP.voteMachine.IsStarted() == true && checkPointId == this.currentChkP.Id {
		// log.Printf("In %s, %s vote %s\n", this.currentChkP.Data(), who, option)
		voteRecordStatus, tallyStatus, newChkPointStatus, fallbackAttempt, outputEvent = this.currentChkP.vote(who, option)
	} else {
		voteRecordStatus = FAILED
	}
	mIdStr := string(this.id)
	if voteRecordStatus == SUCCEED {
		ev.EmitVoteRecorded(mIdStr, who)
	} else if voteRecordStatus == FAILED {
		ev.EmitVoteFailToRecord(mIdStr, who)
	}
	if tallyStatus == SUCCEED {
		ev.EmitTallySucceed(mIdStr, lastChkPointId)
	} else if tallyStatus == FAILED {
		ev.EmitTallyFailed(mIdStr, lastChkPointId)
	}
	if fallbackAttempt == true {
		ev.EmitFallbackAttempt(mIdStr, lastChkPointId)
	}
	if newChkPointStatus == SUCCEED {
		ev.EmitCheckPointStarted(mIdStr, lastChkPointId, this.currentChkP.Id)
	} else if newChkPointStatus == FAILED {
		_, selectedOption := this.currentChkP.voteMachine.GetTallyResult()
		ev.EmitCheckPointFailToStart(mIdStr, lastChkPointId, selectedOption)
	}
	if outputEvent != nil {
		ev.Emit(outputEvent.Id)
	}
	return voteRecordStatus, tallyStatus, newChkPointStatus, fallbackAttempt, outputEvent
}

/**
* Function Beat
* Run at every block
 */
func (this *Mission) BeatAtNewBlock() (ExecutionStatus, bool, ExecutionStatus, *event.Event) {
	ev := event.EmitPredefinedEvent(this.blockchain.GetEventManager())
	var outputEvent *event.Event = nil
	var tallyStatus = DIDNOTSTART
	var newChkPoinStatus = DIDNOTSTART
	var fallbackAttempt bool = false
	var nodeStarted = DIDNOTSTART
	if this.isActive == false {
		return tallyStatus, fallbackAttempt, nodeStarted, outputEvent
	}
	lastChkPointId := this.currentChkP.Id
	mIdStr := string(this.id)
	if this.currentChkP.voteMachine.ShouldTally() == true {
		tallyStatus, newChkPoinStatus, outputEvent = tally(this, this.currentChkP.voteMachine)
		if tallyStatus == SUCCEED {
			ev.EmitTallySucceed(mIdStr, lastChkPointId)
		} else {
			ev.EmitTallyFailed(mIdStr, lastChkPointId)
		}
		if newChkPoinStatus == SUCCEED {
			ev.EmitCheckPointStarted(mIdStr, lastChkPointId, this.currentChkP.Id)
		} else if newChkPoinStatus == FAILED {
			_, selectedOption := this.currentChkP.voteMachine.GetTallyResult()
			ev.EmitCheckPointFailToStart(mIdStr, lastChkPointId, selectedOption)
		}
	} else {
		fallbackAttempt, nodeStarted, outputEvent = fallback(this, this.currentChkP.voteMachine, this.currentChkP.FallbackId)
		if fallbackAttempt == true {
			ev.EmitFallbackAttempt(mIdStr, lastChkPointId)
		}
		if nodeStarted == FAILED {
			_, selectedOption := this.currentChkP.voteMachine.GetTallyResult()
			ev.EmitCheckPointFailToStart(mIdStr, lastChkPointId, selectedOption)
		} else if nodeStarted == SUCCEED {
			ev.EmitCheckPointStarted(mIdStr, lastChkPointId, this.currentChkP.Id)
		}
	}
	return tallyStatus, fallbackAttempt, nodeStarted, outputEvent
}

/**
* Reveal the vote content then tally
 */
func (this *Mission) Reveal(priK []byte) {
	this.currentChkP.voteMachine.Reveal(priK)
}

/**
* Marshal the Mission data
 */
func (this *Mission) marshal() MissionData {
	// go through the tree to take all checkpoint
	return MissionData{}
}

func (this *Mission) unmarshalCheckPoint(data CheckPointData) *CheckPoint {
	return &CheckPoint{}
}
