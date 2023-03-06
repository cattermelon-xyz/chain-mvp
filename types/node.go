package types

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/hectagon-finance/chain-mvp/third_party/tree"
)

const NoFallbackOption = math.MaxUint64 - 1
const EndOfMission = math.MaxUint64 - 2

type Node struct {
	Title       string
	Description string
	FallbackId  uint64
	children    []*Node
	voteMachine VotingMachine
}

// return something that is printable
func (n *Node) Data() interface{} {
	return n.Title
}

// cannot return n.children directly.
// https://github.com/golang/go/wiki/InterfaceSlice
func (n *Node) Children() (c []tree.Node) {
	for _, child := range n.children {
		c = append(c, tree.Node(child))
	}
	return
}

func CreateEmptyNode(title string, desc string, b VotingMachine) *Node {
	node := Node{
		Title:       title,
		Description: desc,
		children:    []*Node{},
		voteMachine: b,
		FallbackId:  NoFallbackOption,
	}
	return &node
}
func CreateNodeWithChildren(name string, desc string, children []*Node, b VotingMachine, fallbackId uint64) *Node {
	node := Node{
		Title:       name,
		Description: desc,
		FallbackId:  fallbackId,
		children:    children,
		voteMachine: b,
	}
	return &node
}
func (this *Node) Attach(child *Node) *Node {
	this.children = append(this.children, child)
	return child
}

/**
* Conversational text the describe the current state of the Node
* including: Title, Description, Options, How voting will conduct
**/
func (this *Node) Print() {
	fmt.Printf("%s\n%s\nVoting Mechanism:\n%s\n", this.Title, this.Description, this.voteMachine.Desc())
	for i := range this.children {
		fmt.Printf("- opt %d: %s\n", i, this.children[i].Title)
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
	if len(this.children) == 0 || this.FallbackId == NoFallbackOption {
		return false
	}
	currentBlockNumber := GetCurrentBlockNumber()
	return this.voteMachine.Start(lastTalliedResult, uint64(len(this.children)), currentBlockNumber, this.FallbackId)
}
func (this *Node) isValidChoice(option []byte) bool {
	if this.voteMachine.IsStarted() == false {
		return false
	}
	return this.voteMachine.ValidateVote(option)
}

/**
* Function vote
* Params: tr *Mission, who string, option []byte
* Returns: voteRecordedSucceed bool, talliedSucceed bool, newNodeStartedSucceed bool
* TODO: what if we want to hide the voter's option from validator?
 */
func (this *Node) vote(tr *Mission, who string, option []byte) (bool, bool, bool) {
	isRecored := this.voteMachine.Record(who, option)
	b, _ := json.Marshal(who)
	if isRecored == true {
		fmt.Println("Vote is recorded")
		_, evId := CreateEvent("VoteRecorded", b)
		Emit(evId)
	} else {
		fmt.Println("Vote record failed")
		_, evId := CreateEvent("VoteRecordFailed", b)
		Emit(evId)
		return false, false, false
	}
	if this.voteMachine.ShouldTally() == true {
		isTallied := this.voteMachine.Tally()
		_, evId := CreateEvent("IsTallied", b)
		Emit(evId)
		if isTallied == true {
			tallyResult, option := this.voteMachine.GetTallyResult()
			fmt.Printf("Tally and the new option is %d\n", option)
			if option != NoOptionMade {
				newNodeStarted, _ := tr.Choose(option, tallyResult)
				fmt.Printf("New node started %t\n", newNodeStarted)
				return true, true, newNodeStarted
			}
		} else {
			fmt.Println("Tally failed")
			return true, false, false
		}
	} else {
		lastTalliedBlockNo := this.voteMachine.GetLastTalliedBlock()
		tallyResult, option := this.voteMachine.GetTallyResult()
		if lastTalliedBlockNo != NeverBeenTallied && option == NoOptionMade {
			// If voteMachine cannot Tally a result, then it should go to the FallbackOption
			newNodeStarted, _ := tr.Choose(this.FallbackId, tallyResult)
			fmt.Printf("New node started %t\n", newNodeStarted)
			return true, true, newNodeStarted
		}
	}
	return true, false, false
}
