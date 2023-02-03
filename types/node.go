package types

import (
	"fmt"

	"github.com/hectagon-finance/chain-mvp/third_party/tree"
)

// voter can choose an option or add few more option thus appending new nodes to the tree
// NOTE: change the name to VotingMachine to reflect the configuration?
type Ballot interface {
	Vote(*Tree, string, int)
	GetName() string
}

const TallyAtVote = -1

// Record to record the data; Tally to take action from the data; TallyAt return the timestamp to active Tally
type VotingMachine interface {
	Record(string, int)
	Tally(*Tree)
	TallyAt() int
}

// NOTE: Record then Tally, instead of Vote!

// system will tally all node in queue

type Tally func(tree *Tree)

type Enforce func(command string) bool

const CHECKPOINT string = "Checkpoint"
const ENFORCER string = "Enforcer"

type Node struct {
	name       string
	enforce    Enforce
	tally      Tally
	children   []*Node
	voteObject Ballot
}

// return something that is printable
func (n *Node) Data() interface{} {
	return n.name
}

// cannot return n.children directly.
// https://github.com/golang/go/wiki/InterfaceSlice
func (n *Node) Children() (c []tree.Node) {
	for _, child := range n.children {
		c = append(c, tree.Node(child))
	}
	return
}

func CreateEmptyNode(name string, b Ballot) *Node {
	node := Node{
		name:       name,
		children:   []*Node{},
		voteObject: b,
	}
	return &node
}
func CreateNodeWithChildren(name string, children []*Node, b Ballot) *Node {
	node := Node{
		name:       name,
		children:   children,
		voteObject: b,
	}
	return &node
}
func (this *Node) Attach(child *Node) *Node {
	this.children = append(this.children, child)
	return child
}
func (this *Node) print() {
	fmt.Printf("%s has following children:\n", this.name)
	for i := range this.children {
		fmt.Printf("- opt %d: %s\n", i, this.children[i].name)
	}
	fmt.Printf("\n")
}
func (this *Node) Get(idx int) *Node {
	if idx < len(this.children) {
		return this.children[idx]
	}
	return nil
}
func (this *Node) isValidChoice(idx int) bool {
	if idx < len(this.children) {
		return true
	}
	return false
}
func (this *Node) vote(tr *Tree, who string, option int) {
	this.voteObject.Vote(tr, who, option)
}

// func (this *Node) vote(idx int) {

// }
