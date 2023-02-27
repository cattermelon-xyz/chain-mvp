package types

import (
	"fmt"

	"github.com/hectagon-finance/chain-mvp/third_party/tree"
)

type Node struct {
	name        string
	children    []*Node
	voteMachine VotingMachine
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
		name:        name,
		children:    []*Node{},
		voteMachine: b,
	}
	return &node
}
func CreateNodeWithChildren(name string, children []*Node, b VotingMachine) *Node {
	node := Node{
		name:        name,
		children:    children,
		voteMachine: b,
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
func (this *Node) Get(idx uint64) *Node {
	if idx < uint64(len(this.children)) {
		return this.children[idx]
	}
	return nil
}
func (this *Node) Start(lastTalliedResult []byte) bool {
	return this.voteMachine.Start(lastTalliedResult, uint64(len(this.children)))
}
func (this *Node) isValidChoice(option interface{}) bool {
	if this.voteMachine.IsStarted() == false {
		return false
	}
	return this.voteMachine.ValidateVote(option)
}

/**
* Function vote
* Params: tr *Initiative, who string, option interface{}
* Returns: voteRecordedSucceed bool, talliedSucceed bool, newNodeStartedSucceed bool
 */
func (this *Node) vote(tr *Initiative, who string, option interface{}) (bool, bool, bool) {
	isRecored := this.voteMachine.Record(who, option)
	if isRecored == true {
		fmt.Println("Vote is recorded")
		_, evId := CreateEvent("VoteRecorded", nil)
		Emit(evId)
	} else {
		fmt.Println("Vote record failed")
		return false, false, false
	}
	if this.voteMachine.TallyAt() == TallyAtVote {
		currentBlockNumber := GetCurrentBlockNumber()
		isTallied := this.voteMachine.Tally(currentBlockNumber)
		if isTallied == true {
			tallyResult, option := this.voteMachine.GetTallyResult()
			fmt.Printf("Tally and the new option is %d\n", option)
			if option != NoOptionMade {
				tr.Choose(option)
				newNodeStarted := tr.Current.Start(tallyResult)
				if newNodeStarted == true {
					fmt.Println("New node started")
					return true, true, true
				} else {
					fmt.Println("New node not started")
					return true, true, false
				}
			}
		} else {
			fmt.Println("Tally falied")
			return true, false, false
		}
	}
	return true, false, false
}
