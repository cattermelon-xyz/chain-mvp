package types

import (
	"errors"
	"fmt"

	"github.com/hectagon-finance/chain-mvp/third_party/tree"
	"github.com/hectagon-finance/chain-mvp/third_party/utils"
)

type Initiative struct {
	id       string
	Title    string
	Fulltext string
	Id       Address
	Owner    Address

	StartNode   *Node
	Current     *Node
	isStarted   bool
	isActivated bool
}

var initiatives = make(map[string]*Initiative)

func CreateInitiative(title string, fulltext string, start *Node) (*Initiative, string) {
	id := utils.RandString(16)
	i := Initiative{
		id:          id,
		Title:       title,
		Fulltext:    fulltext,
		StartNode:   start,
		Current:     nil,
		isStarted:   false,
		isActivated: false,
	}
	initiatives[id] = &i
	return &i, id
}

func GetInitiative(id string) (*Initiative, error) {
	if i, ok := initiatives[id]; ok == true {
		return i, nil
	}
	return nil, errors.New(id + " not found")
}

// TODO: is it safe to do this? should we check all the nodes and events (observer)?
func DeleteInitivate(id string) bool {
	if _, ok := initiatives[id]; ok {
		delete(initiatives, id)
		return true
	}
	return false
}

// func (this *Initiative) edit(d Initiative) bool {
// 	return false
// }

func (this *Initiative) Start() {
	if this.isStarted == false {
		this.isStarted = true
		this.isActivated = true
		this.Current = this.StartNode
		this.Current.Start(nil)
	}
}

func (this *Initiative) Stop() {
	if this.isStarted == true {
		this.isStarted = false
		this.isActivated = false
	}
}

func (this *Initiative) Pause() {
	if this.isActivated == true && this.isStarted == true {
		this.isActivated = false
	}
}

func (this *Initiative) Resume() (bool, error) {
	if this.isStarted == true {
		this.isActivated = true
		return true, nil
	}
	return false, errors.New(this.id + " is stopped, can not start again")
}

func (this *Initiative) PrintFromStart() {
	tree.Print(this.StartNode)
}
func (this *Initiative) PrintFromCurrent() {
	if this.isStarted != true {
		tree.Print(this.StartNode)
	} else {
		tree.Print(this.Current)
	}
}

/**
* TODO: Beside moving to the NextNode, should init something in the nextNode with result from the last Node
**/
func (this *Initiative) Choose(idx uint64) {
	nextNode := this.Current.Get(idx)
	if nextNode == nil {
		fmt.Println(idx, " out of bound, no move")
	}
	if nextNode != nil {
		fmt.Printf("from %s choose: %d got %s\n", this.Current.name, idx, nextNode.name)
		this.Current = nextNode
		this.PrintFromCurrent()
	}
	// emit Event
}
func (this *Initiative) IsValidChoice(option interface{}) bool {
	return this.Current.isValidChoice(option)
}

/**
* Function Vote
* Params: option interface{}, who string
* Returns: voteRecordedSucceed bool, talliedSucceed bool, newNodeStartedSucceed bool
 */
func (this *Initiative) Vote(option interface{}, who string) (bool, bool, bool) {
	if !this.isActivated {
		return false, false, false
	}
	if this.IsValidChoice(option) {
		fmt.Printf("In %s, with %s, %s vote %s\n", this.Current.Data(), this.Current.voteMachine.Desc(), who, option)
		return this.Current.vote(this, who, option)
	}
	return false, false, false
}
