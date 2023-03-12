package types

import (
	"errors"
	"log"
	"strconv"

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
	id       string
	Title    string
	Fulltext string
	Id       Address
	Owner    Address

	EventManager event.EventManager
	startChkP    *CheckPoint
	currentChkP  *CheckPoint
	isStarted    bool
	isActive     bool
}

var Missions = make([]*Mission, 0)

func CreateMission(title string, fulltext string, start *CheckPoint) (*Mission, string) {
	id := utils.RandString(16)
	i := Mission{
		id:           id,
		Title:        title,
		Fulltext:     fulltext,
		startChkP:    start,
		currentChkP:  nil,
		isStarted:    false,
		isActive:     false,
		EventManager: event.GetEventManager(),
	}
	Missions = append(Missions, &i)
	return &i, id
}

func GetMission(id string) (*Mission, error) {
	for _, m := range Missions {
		if m.id == id {
			return m, nil
		}
	}
	return nil, errors.New(id + " not found")
}

// TODO: is it safe to do this? should we check all the nodes and events (observer)?
func DeleteMission(id string) bool {
	found := -1
	for idx, m := range Missions {
		if m.id == id {
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
func (this *Mission) SetStartChkP(chkP *CheckPoint) {
	this.startChkP = chkP
}
func (this *Mission) SetCurrentChkP(chkP *CheckPoint) {
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
	if this.isStarted == false {
		nodeStarted := this.startChkP.start(nil)
		if nodeStarted == false {
			log.Fatal("Mission cannot start")
		} else {
			this.isStarted = true
			this.isActive = true
			this.SetCurrentChkP(this.startChkP)
			this.EventManager.EmitMissionStarted(this.id)
		}
	}
	return this.isStarted
}

func (this *Mission) Stop() {
	if this.isStarted == true {
		this.isStarted = false
		this.isActive = false
		this.EventManager.EmitMissionStopped(this.id)
	}
}

func (this *Mission) Pause() {
	if this.isActive == true && this.isStarted == true {
		this.isActive = false
		this.EventManager.EmitMissionPaused(this.id)
	}
}

func (this *Mission) Resume() (bool, error) {
	if this.isStarted == true {
		this.isActive = true
		this.EventManager.EmitMissionResumed(this.id)
		return true, nil
	}
	return false, errors.New(this.id + " is stopped, can not start again")
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

func (this *Mission) Choose(idx uint64, tallyResult []byte) (bool, error) {
	nextChkP := this.currentChkP.Get(idx)
	started := false
	var err error = nil
	if nextChkP == nil {
		msg := strconv.FormatUint(idx, 10) + " out of bound, no move"
		// log.Fatal(msg)
		log.Println(msg)
		err = errors.New(msg)
	} else {
		// log.Printf("from %s choose: %d got %s\n", this.CurrentChkP.Title, idx, nextChkP.Title)
		started = nextChkP.start(tallyResult)
		if started == true {
			this.currentChkP = nextChkP
		}
	}
	return started, err
}
func (this *Mission) IsValidChoice(option []byte) bool {
	return this.currentChkP.isValidChoice(option)
}

/**
* Function Vote
* Params: option []byte, who string
* Returns: voteRecordedSucceed ExecutionStatus, talliedSucceed ExecutionStatus, newChkPointStartedSucceed ExecutionStatus, fallbackAttempt bool
 */
func (this *Mission) Vote(option []byte, who string) (ExecutionStatus, ExecutionStatus, ExecutionStatus, bool) {
	ev := event.EmitPredefinedEvent(this.EventManager)
	voteRecordStatus := DIDNOTSTART
	tallyStatus := DIDNOTSTART
	newChkPointStatus := DIDNOTSTART
	fallbackAttempt := false
	if this.isActive == false {
		return DIDNOTSTART, DIDNOTSTART, DIDNOTSTART, false
	}
	lastChkPointId := this.currentChkP.Id
	if this.IsValidChoice(option) == true && this.currentChkP.voteMachine.IsStarted() == true {
		log.Printf("In %s, %s vote %s\n", this.currentChkP.Data(), who, option)
		voteRecordStatus, tallyStatus, newChkPointStatus, fallbackAttempt = this.currentChkP.vote(this, who, option)
	} else {
		voteRecordStatus = FAILED
	}
	if voteRecordStatus == SUCCEED {
		ev.EmitVoteRecorded(this.id, who)
	} else if voteRecordStatus == FAILED {
		ev.EmitVoteFailToRecord(this.id, who)
	}
	if tallyStatus == SUCCEED {
		ev.EmitTallySucceed(this.id, lastChkPointId)
	} else if tallyStatus == FAILED {
		ev.EmitTallyFailed(this.id, lastChkPointId)
	}
	if fallbackAttempt == true {
		ev.EmitFallbackAttempt(this.id, lastChkPointId)
	}
	if newChkPointStatus == SUCCEED {
		ev.EmitCheckPointStarted(this.id, lastChkPointId, this.currentChkP.Id)
	} else if newChkPointStatus == FAILED {
		_, selectedOption := this.currentChkP.voteMachine.GetTallyResult()
		ev.EmitCheckPointFailToStart(this.id, lastChkPointId, selectedOption)
	}
	return voteRecordStatus, tallyStatus, newChkPointStatus, fallbackAttempt
}

/**
* Function Tally
* Run at every block
 */
func (this *Mission) TallyAtNewBlock() {
	ev := event.EmitPredefinedEvent(this.EventManager)
	if this.isActive == false {
		return
	}
	lastChkPointId := this.currentChkP.Id
	if this.currentChkP.voteMachine.ShouldTally() == true {
		tallyStatus, newChkPoinStatus := tally(this, this.currentChkP.voteMachine)
		if tallyStatus == SUCCEED {
			ev.EmitTallySucceed(this.id, lastChkPointId)
		} else {
			ev.EmitTallyFailed(this.id, lastChkPointId)
		}
		if newChkPoinStatus == SUCCEED {
			ev.EmitCheckPointStarted(this.id, lastChkPointId, this.currentChkP.Id)
		} else if newChkPoinStatus == FAILED {
			_, selectedOption := this.currentChkP.voteMachine.GetTallyResult()
			ev.EmitCheckPointFailToStart(this.id, lastChkPointId, selectedOption)
		}
	} else {
		fallbackAttempt, nodeStarted := fallback(this, this.currentChkP.voteMachine, this.currentChkP.FallbackId)
		if fallbackAttempt == true {
			ev.EmitFallbackAttempt(this.id, lastChkPointId)
		}
		if nodeStarted == FAILED {
			_, selectedOption := this.currentChkP.voteMachine.GetTallyResult()
			ev.EmitCheckPointFailToStart(this.id, lastChkPointId, selectedOption)
		} else if nodeStarted == SUCCEED {
			ev.EmitCheckPointStarted(this.id, lastChkPointId, this.currentChkP.Id)
		}
	}
}

/**
* Reveal the vote content then tally
 */
func (this *Mission) Reveal(priK []byte) {
	this.currentChkP.voteMachine.Reveal(priK)
}
