package types

// voter can choose an option or add few more option thus appending new nodes to the tree
// NOTE: change the name to VotingMachine to reflect the configuration?
type Ballot interface {
	Vote(*Initiative, string, int)
	GetName() string
}

const TallyAtVote = -1

// Record to record the data; Tally to take action from the data; TallyAt return the timestamp to active Tally
type VotingMachine interface {
	Record(string, int)
	Tally(*Initiative)
	TallyAt() int
	IsTallied() bool
}
