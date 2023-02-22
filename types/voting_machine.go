package types

// voter can choose an option or add few more option thus appending new nodes to the tree
// NOTE: change the name to VotingMachine to reflect the configuration?
type Ballot interface {
	Vote(*Initiative, string, int)
	GetName() string
}

const TallyAtVote = -1
const NoTallyResult = -1

// Record to record the data; Tally to take action from the data; TallyAt return the timestamp to active Tally
type VotingMachine interface {
	// Describe the rule of the vote
	Desc() string
	// Record the data: who <string> choose option <int>
	Record(string, int)
	// Calculate the voting power of the vote
	VotingPower(string, int) int
	// Cost of the Vote
	Cost(string, int) int
	// Tally the vote
	Tally()
	// When to tally the vote; if TallyAt() != -1 then it can only tally ONCE
	TallyAt() int
	// Return the Tally result, return NoTallyResult if no option is made
	GetTallyResult() int
}
