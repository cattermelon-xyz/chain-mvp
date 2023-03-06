package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/hectagon-finance/chain-mvp/third_party/tree"
	"github.com/hectagon-finance/chain-mvp/third_party/utils"
)

type Mission struct {
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

var Missions = make(map[string]*Mission)

func CreateMission(title string, fulltext string, start *Node) (*Mission, string) {
	id := utils.RandString(16)
	i := Mission{
		id:          id,
		Title:       title,
		Fulltext:    fulltext,
		StartNode:   start,
		Current:     nil,
		isStarted:   false,
		isActivated: false,
	}
	Missions[id] = &i
	return &i, id
}

func GetMission(id string) (*Mission, error) {
	if i, ok := Missions[id]; ok == true {
		return i, nil
	}
	return nil, errors.New(id + " not found")
}

// TODO: is it safe to do this? should we check all the nodes and events (observer)?
func DeleteInitivate(id string) bool {
	if _, ok := Missions[id]; ok {
		delete(Missions, id)
		return true
	}
	return false
}

// func (this *Mission) edit(d Mission) bool {
// 	return false
// }

func (this *Mission) Start() bool {
	if this.isStarted == false {
		nodeStarted := this.StartNode.Start(nil)
		if nodeStarted == false {
			fmt.Println("Mission cannot start")
		} else {
			this.isStarted = true
			this.isActivated = true
			this.Current = this.StartNode
			fmt.Println("Mission started successfully")
		}
	}
	return this.isStarted
}

func (this *Mission) Stop() {
	if this.isStarted == true {
		this.isStarted = false
		this.isActivated = false
	}
}

func (this *Mission) Pause() {
	if this.isActivated == true && this.isStarted == true {
		this.isActivated = false
	}
}

func (this *Mission) Resume() (bool, error) {
	if this.isStarted == true {
		this.isActivated = true
		return true, nil
	}
	return false, errors.New(this.id + " is stopped, can not start again")
}

func (this *Mission) PrintFromStart() {
	tree.Print(this.StartNode)
}
func (this *Mission) PrintFromCurrent() {
	if this.isStarted != true {
		tree.Print(this.StartNode)
	} else {
		tree.Print(this.Current)
	}
}

func (this *Mission) Choose(idx uint64, tallyResult []byte) (bool, error) {
	nextNode := this.Current.Get(idx)
	started := false
	var err error = nil
	if nextNode == nil {
		msg := strconv.FormatUint(idx, 64) + " out of bound, no move"
		log.Fatal(msg)
		err = errors.New(msg)
	} else {
		fmt.Printf("from %s choose: %d got %s\n", this.Current.Title, idx, nextNode.Title)
		this.Current = nextNode
		this.PrintFromCurrent()
		started = nextNode.Start(tallyResult)

		args := []string{this.id, nextNode.Title}
		b, _ := json.Marshal(args)
		_, evid := CreateEvent("NextCheckPoint", b)
		Emit(evid)
	}
	return started, err
}
func (this *Mission) IsValidChoice(option []byte) bool {
	return this.Current.isValidChoice(option)
}

/**
* Function Vote
* Params: option []byte, who string
* Returns: voteRecordedSucceed bool, talliedSucceed bool, newNodeStartedSucceed bool
 */
func (this *Mission) Vote(option []byte, who string) (bool, bool, bool) {
	if !this.isActivated {
		return false, false, false
	}
	if this.IsValidChoice(option) {
		fmt.Printf("In %s, with %s, %s vote %s\n", this.Current.Data(), this.Current.voteMachine.Desc(), who, option)
		return this.Current.vote(this, who, option)
	}
	return false, false, false
}
