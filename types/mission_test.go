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
	blkc := &MockBlockchain{}
	ev := MockEventManager{}
	blkc.SetEventManager(&ev)
	missionA, _ := types.CreateMission("missionA", "test", blkc)
	chkP := missionA.CreateEmptyCheckPoint("test", "test desc", &alwaysStartMachine)
	chkP.FallbackId = 1
	chkp1 := missionA.CreateEmptyCheckPoint("test2", "test desc 2", &alwaysStartMachine)
	chkP.Attach(chkp1.Id)

	missionA.SetStartChkP(chkP.Id)
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
	ev := MockEventManager{}
	blkc := &MockBlockchain{}
	blkc.SetEventManager(&ev)
	machine := MockVoteMachine{
		Started: true,
	}
	missionA, _ := types.CreateMission("missionA", "test", blkc)
	chkP := missionA.CreateEmptyCheckPoint("test", "test desc", &machine)
	chkP.FallbackId = 0
	chkP2 := missionA.CreateEmptyCheckPoint("test2", "test desc 2", &machine)
	chkP.Attach(chkP2.Id)
	missionA.SetStartChkP(chkP.Id)

	var recordState, tallyState, newNodeState types.ExecutionStatus
	var fallbackAttemp bool
	recordState, tallyState, newNodeState, fallbackAttemp, _ = missionA.Vote([]byte{0}, "", "")
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
	machine = MockVoteMachine{
		Started:          true,
		VoteValid:        false,
		ShouldTallyState: false,
	}
	blkc.SetCurrentBlockNumber(100)
	missionB, _ := types.CreateMission("missionB", "desc", blkc)
	chkP = missionB.CreateCheckPoinWithChildren("[ChkP] test name", "[ChkP] test desc",
		[]*types.CheckPoint{{}}, &machine, 2, uint64(1000), uint64(1000))
	missionB.SetStartChkP(chkP.Id)
	missionB.Start()
	assert.Equal(t, ev.UnQueue(), event.MissionStarted, "Mission should Started")
	recordState, _, _, _, _ = missionB.Vote([]byte{0}, "x", chkP.Id)
	assert.Equal(t, recordState, types.FAILED, "Record should fail")
	assert.Equal(t, ev.UnQueue(), event.VoteFailToRecord, "Should emit VoteFailToRecord")
	machine.VoteValid = true
	machine.VoteRecordSucceed = false
	recordState, _, _, _, _ = missionB.Vote([]byte{0}, "x", chkP.Id)
	assert.Equal(t, recordState, types.FAILED, "Record should fail")
	assert.Equal(t, ev.UnQueue(), event.VoteFailToRecord, "Should emit VoteFailToRecord")
	// vote failed to record: mismatch checkPointId
	machine.VoteRecordSucceed = true
	recordState, _, _, _, _ = missionB.Vote([]byte{0}, "x", "Incorrect CheckPointId")
	assert.Equal(t, recordState, types.FAILED, "Record should fail")
	assert.Equal(t, ev.UnQueue(), event.VoteFailToRecord, "Should emit VoteFailToRecord")
	// vote recorded, no tally
	machine.VoteRecordSucceed = true
	recordState, _, _, _, _ = missionB.Vote([]byte{0}, "x", chkP.Id)
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
	missionA, _ := types.CreateMission("missionA", "desc missionA", blkc)
	chkP2 := missionA.CreateEmptyCheckPoint("CheckPoint2", "test desc", &machine2)
	chkP2.FallbackId = 0
	chkP3 := missionA.CreateEmptyCheckPoint("CheckPoint3", "test desc", &machine2)
	chkP2.Attach(chkP3.Id)
	chkP := missionA.CreateCheckPoinWithChildren("[ChkP] test name", "[ChkP] test desc",
		[]*types.CheckPoint{chkP2}, &machine, 0, uint64(1000), uint64(1000))
	missionA.SetStartChkP(chkP.Id)
	missionA.Start()
	assert.Equal(t, ev.UnQueue(), event.MissionStarted, "Mission should Started")
	var recordState, tallyState, newNodeState types.ExecutionStatus
	var fallbackAttempt bool
	// vote recorded, fallback but new node fail to start
	machine.VoteRecordSucceed = true
	machine.ShouldTallyState = false
	machine2.Started = false
	blkc.SetCurrentBlockNumber(1110) // currentBlock > last vote & last tally
	recordState, tallyState, newNodeState, fallbackAttempt, _ = missionA.Vote([]byte{0}, "x", chkP.Id)
	assert.Equal(t, recordState, types.DIDNOTSTART, "Should NOT be DIDNOTSTART")
	assert.Equal(t, fallbackAttempt, true, "fallbackAttempt should be TRUE")
	assert.Equal(t, tallyState, types.DIDNOTSTART, "tally should be DIDNOTSTART")
	assert.Equal(t, newNodeState, types.FAILED, "New node failed to start")
	assert.Equal(t, ev.UnQueue(), event.FallbackAttempt, "Should emit FallbackAttempt")
	assert.Equal(t, ev.UnQueue(), event.CheckPointFailToStart, "Should emit CheckPointFailToStart")
	// vote recorded, fallback, new node started
	machine2.Started = true
	chkP = missionA.CreateCheckPoinWithChildren("[ChkP] test name", "[ChkP] test desc",
		[]*types.CheckPoint{chkP2}, &machine, 0, uint64(1000), uint64(1000)) // reset
	missionA.SetCurrentChkP(chkP.Id)
	recordState, tallyState, newNodeState, fallbackAttempt, _ = missionA.Vote([]byte{0}, "x", chkP.Id)
	assert.Equal(t, recordState, types.DIDNOTSTART, "Should be DIDNOTSTART")
	assert.Equal(t, fallbackAttempt, true, "fallbackAttempt should be TRUE")
	assert.Equal(t, tallyState, types.DIDNOTSTART, "tally should be DIDNOTSTART")
	assert.Equal(t, newNodeState, types.SUCCEED, "New node should start successfully")
	assert.Equal(t, ev.UnQueue(), event.FallbackAttempt, "Should emit FallbackAttempt")
	assert.Equal(t, ev.UnQueue(), event.CheckPointStarted, "Should emit CheckPointStarted")
	// vote recorded, fallback, new node is an Output
	ev.Clear()
	missionB, _ := types.CreateMission("missionB", "fulltext missionB", blkc)
	e, _ := ev.CreateEvent("Output1", []byte{0})
	output1 := missionB.CreateOutput("output1", "output1 desc", e)
	chkP4 := missionB.CreateCheckPoinWithChildren("test", "test desc", []*types.CheckPoint{
		output1,
	}, &machine, 0, 0, 0)
	missionB.SetStartChkP(chkP4.Id)
	missionB.Start()
	assert.Equal(t, ev.UnQueue(), event.MissionStarted, "Should emit MissionStarted")
	recordState, tallyState, newNodeState, fallbackAttempt, _ = missionB.Vote([]byte{0}, "x", chkP4.Id)
	assert.Equal(t, recordState, types.DIDNOTSTART, "Should be DIDNOTSTART")
	assert.Equal(t, fallbackAttempt, true, "fallbackAttempt should be TRUE")
	assert.Equal(t, tallyState, types.DIDNOTSTART, "tally should be DIDNOTSTART")
	assert.Equal(t, newNodeState, types.SUCCEED, "New node should start successfully")
	assert.Equal(t, ev.UnQueue(), event.MissionStopped, "Should emit MissionStoped")
	assert.Equal(t, ev.UnQueue(), event.FallbackAttempt, "Should emit FallbackAttempt")
	assert.Equal(t, ev.UnQueue(), event.CheckPointStarted, "Should emit CheckPointStarted")
	assert.Equal(t, ev.UnQueue(), "Output1", "Should emit Output1")
}

func TestMissionVoteOnTally(t *testing.T) {
	// vote recorded, no tally
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
	missionA, _ := types.CreateMission("missionA", "desc missionA", blkc)
	chkP2 := missionA.CreateEmptyCheckPoint("CheckPoint2", "test desc", &machine2)
	chkP2.FallbackId = 0
	chkP3 := missionA.CreateEmptyCheckPoint("CheckPoint3", "test desc", &machine2)
	chkP2.Attach(chkP3.Id)
	chkP := missionA.CreateCheckPoinWithChildren("[ChkP] test name", "[ChkP] test desc",
		[]*types.CheckPoint{chkP2}, &machine, 0, uint64(1000), uint64(1000))
	missionA.SetStartChkP(chkP.Id)
	missionA.Start()
	// vote recorded, tally failed
	machine.VoteRecordSucceed = true
	machine.ShouldTallyState = true
	machine.TallyExecutionState = false
	recordState, tallyState, newNodeState, fallbackAttempt, _ := missionA.Vote([]byte{0}, "x", chkP.Id)
	assert.Equal(t, recordState, types.SUCCEED, "Vote should be succeed")
	assert.Equal(t, tallyState, types.FAILED, "Tally should failed")
	assert.Equal(t, newNodeState, types.DIDNOTSTART, "newNodeState should be DIDNOTSTARTED")
	assert.Equal(t, fallbackAttempt, false, "No fallback attempt")
	assert.Equal(t, ev.UnQueue(), event.MissionStarted, "Mission should be started")
	assert.Equal(t, ev.UnQueue(), event.VoteRecorded, "Vote should be recorded")
	assert.Equal(t, ev.UnQueue(), event.TallyFailed, "Tally should be failed")
	assert.Equal(t, ev.UnQueue(), "", "Event queue should be empty")
	// vote recorded, tally succeed, no new node
	machine.TallyExecutionState = true
	machine.OptionMade = types.NoOptionMade
	recordState, tallyState, newNodeState, fallbackAttempt, _ = missionA.Vote([]byte{0}, "x", chkP.Id)
	assert.Equal(t, recordState, types.SUCCEED, "Vote should be succeed")
	assert.Equal(t, tallyState, types.SUCCEED, "Tally should succeed")
	assert.Equal(t, newNodeState, types.DIDNOTSTART, "newNodeState should be DIDNOTSTARTED")
	assert.Equal(t, fallbackAttempt, false, "No fallback attempt")
	assert.Equal(t, ev.UnQueue(), event.VoteRecorded, "Vote should be recorded")
	assert.Equal(t, ev.UnQueue(), event.TallySucceed, "Tally should be succeed")
	assert.Equal(t, ev.UnQueue(), "", "Event queue should be empty")
	// vote recorded, tally succeed, new node failed
	machine.TallyExecutionState = true
	machine.OptionMade = 0
	machine2.Started = false
	recordState, tallyState, newNodeState, fallbackAttempt, _ = missionA.Vote([]byte{0}, "x", chkP.Id)
	assert.Equal(t, recordState, types.SUCCEED, "Vote should be succeed")
	assert.Equal(t, tallyState, types.SUCCEED, "Tally should succeed")
	assert.Equal(t, newNodeState, types.FAILED, "newNodeState should be FAILED")
	assert.Equal(t, fallbackAttempt, false, "No fallback attempt")
	assert.Equal(t, ev.UnQueue(), event.VoteRecorded, "Vote should be recorded")
	assert.Equal(t, ev.UnQueue(), event.TallySucceed, "Tally should be succeed")
	assert.Equal(t, ev.UnQueue(), event.CheckPointFailToStart, "CheckPoint should fail to start")
	assert.Equal(t, ev.UnQueue(), "", "Event queue should be empty")
	// vote recorded, tally succeed, new node started
	machine.TallyExecutionState = true
	machine.OptionMade = 0
	machine2.Started = true
	recordState, tallyState, newNodeState, fallbackAttempt, _ = missionA.Vote([]byte{0}, "x", chkP.Id)
	assert.Equal(t, recordState, types.SUCCEED, "Vote should be succeed")
	assert.Equal(t, tallyState, types.SUCCEED, "Tally should succeed")
	assert.Equal(t, newNodeState, types.SUCCEED, "newNodeState should be FAILED")
	assert.Equal(t, fallbackAttempt, false, "No fallback attempt")
	assert.Equal(t, ev.UnQueue(), event.VoteRecorded, "Vote should be recorded")
	assert.Equal(t, ev.UnQueue(), event.TallySucceed, "Tally should be succeed")
	assert.Equal(t, ev.UnQueue(), event.CheckPointStarted, "CheckPoint should be started")
	assert.Equal(t, ev.UnQueue(), "", "Event queue should be empty")
}

func TestBeatAtNewBlock(t *testing.T) {
	// no tally start

	// tally failed

	// tally succeed, no new node to start

	// tally succeed, new node fail to start

	// tally succeed, new node started
}
