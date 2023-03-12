package mock

import "github.com/hectagon-finance/chain-mvp/types"

type MockVoteMachine struct {
	Started             bool
	VoteValid           bool
	VoteRecordSucceed   bool
	ShouldTallyState    bool
	TallyExecutionState bool
	OptionMade          uint64
	LastTalliedBlock    uint64
	TallyResult         []byte
	RevealState         bool
	VotingPower         uint64
	Cost                uint64
}

func (this *MockVoteMachine) Desc() string {
	return ""
}
func (this *MockVoteMachine) Start(tallyResult []byte, noOfOptions uint64, fallbackOption uint64) bool {
	return this.Started
}
func (this *MockVoteMachine) IsStarted() bool {
	return this.Started
}
func (this *MockVoteMachine) ValidateVote(option []byte) bool {
	return this.VoteValid
}
func (this *MockVoteMachine) GetVotingPower(who string, option []byte) uint64 {
	return this.VotingPower
}
func (this *MockVoteMachine) GetCost(who string, option []byte) uint64 {
	return this.Cost
}
func (this *MockVoteMachine) Record(who string, option []byte) bool {
	return this.VoteRecordSucceed
}
func (this *MockVoteMachine) ShouldTally() bool {
	return this.ShouldTallyState
}
func (this *MockVoteMachine) Tally() (talliedSucceed bool, lastTalliedBlock uint64, tallyResult []byte, selectedOption uint64) {
	return this.TallyExecutionState, this.LastTalliedBlock, this.TallyResult, this.OptionMade
}
func (this *MockVoteMachine) GetLastTalliedBlock() uint64 {
	return this.LastTalliedBlock
}
func (this *MockVoteMachine) GetTallyResult() (tallyResult []byte, selectedOption uint64) {
	return this.TallyResult, types.NoOptionMade
}
func (this *MockVoteMachine) Reveal(key []byte) bool {
	return this.RevealState
}
