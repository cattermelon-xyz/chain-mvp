package types

import (
	"fmt"

	"github.com/hectagon-finance/chain-mvp/third_party/tree"
)

type Node struct {
	name           string
	children       []*Node
	voteObject     VotingMachine
	activatedEvent Event
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

func CreateEmptyNode(name string, b VotingMachine) *Node {
	node := Node{
		name:       name,
		children:   []*Node{},
		voteObject: b,
		activatedEvent: Event{
			Name: "NodeActivated",
			Args: []string{name},
		},
	}
	return &node
}
func CreateNodeWithChildren(name string, children []*Node, b VotingMachine) *Node {
	node := Node{
		name:       name,
		children:   children,
		voteObject: b,
		activatedEvent: Event{
			Name: "NodeActivated",
			Args: []string{name},
		},
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
func (this *Node) vote(tr *Initiative, who string, option int) {
	this.voteObject.Record(who, option)
	if this.voteObject.TallyAt() == -1 {
		this.voteObject.Tally()
		tallyResult := this.voteObject.GetTallyResult()
		fmt.Printf("Tally and the result %d\n", tallyResult)
		if tallyResult != NoTallyResult {
			tr.Choose(tallyResult)
		}
	}
}
