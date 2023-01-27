package checkpoint

import (
	"github.com/hectagon-finance/chain-mvp/pkg/enforcer"
	"github.com/hectagon-finance/chain-mvp/pkg/net"
)

// each CheckPoint should point to the next through an mutable option list
type CheckPoint struct {
	Title    string
	Fulltext string
	Options  *[]VoteOutput
	VoteFunc *VoteFunc
}

// voteFunc should calculate the voting power and yield the corresponding output
// voteFunc should define the function and input
// Q: how to match voteFunc and options?
type VoteFunc struct {
	Network net.Network
	Address string
}

type VoteOutput struct {
	Next     *CheckPoint
	Enforcer *enforcer.Enforcer
}
