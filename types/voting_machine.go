package types

import "math"

// voter can choose an option or add few more option thus appending new nodes to the tree
// NOTE: change the name to VotingMachine to reflect the configuration?
type Ballot interface {
	Vote(*Mission, string, int)
	GetName() string
}

const TallyAtVote = math.MaxUint64
const NoOptionMade = math.MaxUint64
const NeverBeenTallied = math.MaxUint64

/*
* Record to record the data; Tally to take action from the data; TallyAt return the timestamp to active Tally
 */
type VotingMachine interface {
	// Describe the rule of the vote
	Desc() string
	// Validate the vote
	ValidateVote(option []byte) bool
	// Record the data: who_string choose option_[]byte
	Record(who string, option []byte) bool
	// Calculate the voting power of the vote
	VotingPower(who string, option []byte) uint64
	// Cost of the Vote
	Cost(who string, option []byte) uint64
	// Tally the vote, return if tally happen successfully
	Tally() bool
	// Return true if VotingMachine beable to tally
	ShouldTally() bool
	// Return the last tallied block
	GetLastTalliedBlock() uint64
	// Return the Tally result, return nil []byte and NoOptionMade code if no option made.
	GetTallyResult() ([]byte, uint64)
	// After this, the machine is ready for vote. Return if Start succeed
	Start(tallyResult []byte, noOfOptions uint64, startedBlock uint64, fallbackOption uint64) bool
	// Return the Readiness of the VotingMachine
	IsStarted() bool
}
