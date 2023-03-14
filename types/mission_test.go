package types_test

import (
	"testing"

	. "github.com/hectagon-finance/chain-mvp/mock"
	"github.com/hectagon-finance/chain-mvp/types"
	"github.com/hectagon-finance/chain-mvp/types/event"
	"github.com/stretchr/testify/assert"
)

func TestMissionStateMachine(t *testing.T) {
	alwaysStartMachine := MockVoteMachine{
		Started: true,
	}
	chkP := types.CreateEmptyCheckPoint("test", "test desc", &alwaysStartMachine, nil)
	chkP.Attach(&types.CheckPoint{})

	ev := MockEventManager{}
	missionA := types.Mission{
		EventManager: &ev, // change EventManager to interface
	}
	missionA.SetStartChkP(chkP)
	missionA.Start()
	assert.Equal(t, missionA.Start(), true, "Mission started should return true")
	assert.Equal(t, ev.UnQueue(), event.MissionStarted, "Event MissionStarted must be dispatched")
	ev.Clear()
	missionA.Pause()
	assert.Equal(t, ev.UnQueue(), event.MissionPaused, "Event MissionPaused must be dispatched")
	missionA.Resume()
	assert.Equal(t, ev.UnQueue(), event.MissionResumed, "Event MissionResumed must be dispatched")
	missionA.Stop()
	resumeState, _ := missionA.Resume()
	assert.Equal(t, resumeState, false, "Stopped Mission cannnot Resume")
	started, _ := missionA.Resume()
	assert.Equal(t, started, false, "Stopped Mission cannnot Start")
}

func TestMissionVoteOnRecord(t *testing.T) {
	machine := MockVoteMachine{
		Started: true,
	}
	chkP := types.CreateEmptyCheckPoint("test", "test desc", &machine, nil)
	chkP.Attach(&types.CheckPoint{})

	ev := MockEventManager{}
	missionA := types.Mission{
		EventManager: &ev, // change EventManager to interface
	}
	missionA.SetStartChkP(chkP)

	var recordState, tallyState, newNodeState types.ExecutionStatus
	var fallbackAttemp bool
	recordState, tallyState, newNodeState, fallbackAttemp = missionA.Vote([]byte{0}, "", "")
	// cannot vote a mission not started, paused, stopped
	assert.Equal(t, recordState, types.DIDNOTSTART, "Record must be DIDNOTSTART")
	assert.Equal(t, tallyState, types.DIDNOTSTART, "Tally must be DIDNOTSTART")
	assert.Equal(t, newNodeState, types.DIDNOTSTART, "NewNode must be DIDNOTSTART")
	assert.Equal(t, fallbackAttemp, false, "fallbackAtempt must be false")
	missionA.Start()
	assert.Equal(t, ev.UnQueue(), event.MissionStarted, "Mission must be Started")
	missionA.Pause()
	assert.Equal(t, ev.UnQueue(), event.MissionPaused, "Mission must be Paused")
	assert.Equal(t, recordState, types.DIDNOTSTART, "Record must be DIDNOTSTART")
	assert.Equal(t, tallyState, types.DIDNOTSTART, "Tally must be DIDNOTSTART")
	assert.Equal(t, newNodeState, types.DIDNOTSTART, "NewNode must be DIDNOTSTART")
	assert.Equal(t, fallbackAttemp, false, "fallbackAtempt must be false")
	missionA.Resume()
	assert.Equal(t, ev.UnQueue(), event.MissionResumed, "Mission must be Resumed")
	missionA.Pause()
	assert.Equal(t, ev.UnQueue(), event.MissionPaused, "Mission must be Paused")
	assert.Equal(t, recordState, types.DIDNOTSTART, "Record must be DIDNOTSTART")
	assert.Equal(t, tallyState, types.DIDNOTSTART, "Tally must be DIDNOTSTART")
	assert.Equal(t, newNodeState, types.DIDNOTSTART, "NewNode must be DIDNOTSTART")
	assert.Equal(t, fallbackAttemp, false, "fallbackAtempt must be false")
	missionA.Stop()
	assert.Equal(t, ev.UnQueue(), event.MissionStopped, "Mission must be Stopped")
	assert.Equal(t, recordState, types.DIDNOTSTART, "Record must be DIDNOTSTART")
	assert.Equal(t, tallyState, types.DIDNOTSTART, "Tally must be DIDNOTSTART")
	assert.Equal(t, newNodeState, types.DIDNOTSTART, "NewNode must be DIDNOTSTART")
	assert.Equal(t, fallbackAttemp, false, "fallbackAtempt must be false")
	// vote failed to record: invalid choice
	ev.Clear()
	machine = MockVoteMachine{
		Started:          true,
		VoteValid:        false,
		ShouldTallyState: false,
	}
	blkc := &MockBlockchain{} // no fallback: currentBlock < lastVote, lastTally
	blkc.SetEventManager(&ev)
	blkc.SetCurrentBlockNumber(100)
	chkP = types.CreateCheckPoinWithChildren("[ChkP] test name", "[ChkP] test desc",
		[]*types.CheckPoint{{}}, &machine, 2, uint64(1000), uint64(1000), blkc)
	missionB := types.Mission{
		EventManager: &ev, // change EventManager to interface
	}
	missionB.SetStartChkP(chkP)
	missionB.Start()
	assert.Equal(t, ev.UnQueue(), event.MissionStarted, "Mission should Started")
	recordState, _, _, _ = missionB.Vote([]byte{0}, "x", chkP.Id)
	assert.Equal(t, recordState, types.FAILED, "Record should fail")
	assert.Equal(t, ev.UnQueue(), event.VoteFailToRecord, "Should emit VoteFailToRecord")
	machine.VoteValid = true
	machine.VoteRecordSucceed = false
	recordState, _, _, _ = missionB.Vote([]byte{0}, "x", chkP.Id)
	assert.Equal(t, recordState, types.FAILED, "Record should fail")
	assert.Equal(t, ev.UnQueue(), event.VoteFailToRecord, "Should emit VoteFailToRecord")
	// vote failed to record: mismatch checkPointId
	machine.VoteRecordSucceed = true
	recordState, _, _, _ = missionB.Vote([]byte{0}, "x", "Incorrect CheckPointId")
	assert.Equal(t, recordState, types.FAILED, "Record should fail")
	assert.Equal(t, ev.UnQueue(), event.VoteFailToRecord, "Should emit VoteFailToRecord")
	// vote recorded, no tally
	machine.VoteRecordSucceed = true
	recordState, _, _, _ = missionB.Vote([]byte{0}, "x", chkP.Id)
	assert.Equal(t, recordState, types.SUCCEED, "Vote succeed")
	assert.Equal(t, ev.UnQueue(), event.VoteRecorded, "Should emit VoteRecorded")
}

func TestMissionVoteOnFallback(t *testing.T) {
	blkc := &MockBlockchain{}
	machine := MockVoteMachine{
		Started:   true,
		VoteValid: true,
	}
	machine2 := MockVoteMachine{
		Started: true,
	}

	ev := MockEventManager{}
	blkc.SetEventManager(&ev)
	missionA := types.Mission{
		EventManager: &ev, // change EventManager to interface
	}
	chkP2 := types.CreateEmptyCheckPoint("CheckPoint2", "test desc", &machine2, nil)
	chkP := *types.CreateCheckPoinWithChildren("[ChkP] test name", "[ChkP] test desc",
		[]*types.CheckPoint{chkP2}, &machine, 0, uint64(1000), uint64(1000), blkc)
	missionA.SetStartChkP(&chkP)
	missionA.Start()
	assert.Equal(t, ev.UnQueue(), event.MissionStarted, "Mission should Started")
	var recordState, tallyState, newNodeState types.ExecutionStatus
	var fallbackAttempt bool
	// vote recorded, fallback but new node fail to start
	machine.VoteRecordSucceed = true
	machine.ShouldTallyState = false
	machine2.Started = false
	blkc.SetCurrentBlockNumber(1110) // currentBlock > last vote & last tally
	recordState, tallyState, newNodeState, fallbackAttempt = missionA.Vote([]byte{0}, "x", chkP.Id)
	assert.Equal(t, recordState, types.DIDNOTSTART, "Should NOT be DIDNOTSTART")
	assert.Equal(t, fallbackAttempt, true, "fallbackAttempt should be TRUE")
	assert.Equal(t, tallyState, types.DIDNOTSTART, "tally should be DIDNOTSTART")
	assert.Equal(t, newNodeState, types.FAILED, "New node failed to start")
	assert.Equal(t, ev.UnQueue(), event.FallbackAttempt, "Should emit FallbackAttempt")
	assert.Equal(t, ev.UnQueue(), event.CheckPointFailToStart, "Should emit CheckPointFailToStart")
	// vote recorder, fallback, new node started
	machine2.Started = true
	chkP = *types.CreateCheckPoinWithChildren("[ChkP] test name", "[ChkP] test desc",
		[]*types.CheckPoint{chkP2}, &machine, 0, uint64(1000), uint64(1000), blkc) // reset
	missionA.SetCurrentChkP(&chkP)
	recordState, tallyState, newNodeState, fallbackAttempt = missionA.Vote([]byte{0}, "x", chkP.Id)
	assert.Equal(t, recordState, types.DIDNOTSTART, "Should NOT be DIDNOTSTART")
	assert.Equal(t, fallbackAttempt, true, "fallbackAttempt should be TRUE")
	assert.Equal(t, tallyState, types.DIDNOTSTART, "tally should be DIDNOTSTART")
	assert.Equal(t, newNodeState, types.SUCCEED, "New node start successfully")
	assert.Equal(t, ev.UnQueue(), event.FallbackAttempt, "Should emit FallbackAttempt")
	assert.Equal(t, ev.UnQueue(), event.CheckPointStarted, "Should emit CheckPointStarted")
}

func TestMissionVoteOnTally(t *testing.T) {
	// vote recorded, no tally
	// vote recorded, tally failed
	// vote recorded, tally succeed, no new node

	// machine.VoteRecordSucceed = true
	// machine.ShouldTallyState = true
	// machine.TallyExecutionState = true
	// machine.OptionMade = types.NoOptionMade
	// recordState, tallyState, newNodeState, _ = missionB.Vote([]byte{0}, "")
	// assert.Equal(t, recordState, types.SUCCEED, "Vote succeed")
	// assert.Equal(t, tallyState, types.SUCCEED, "Tally must be Succeed")
	// assert.Equal(t, newNodeState, types.DIDNOTSTART, "Should be DIDNOTSTART")
	// assert.Equal(t, ev.UnQueue(), event.VoteRecorded, "Should emit VoteRecorded")
	// assert.Equal(t, ev.UnQueue(), event.TallySucceed, "Should emit TallySucceed")

	// vote recorded, tally succeed, new node failed
	// vote recorded, tally succeed, new node started
}

func TestMissionTallyAtBlock(t *testing.T) {
	// no tally start

	// tally failed

	// tally succeed, no new node to start

	// tally succeed, new node fail to start

	// tally succeed, new node started
}
