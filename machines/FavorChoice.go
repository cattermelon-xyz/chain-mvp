package machines

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/go-git/go-git/v5/utils/binary"
	"github.com/hectagon-finance/chain-mvp/types"
)

type FavorChoice struct {
	Threshold          uint64
	voted              map[uint64]uint64
	records            map[string]bool
	noOfOption         uint64
	selectedOption     uint64
	lastVotedOpt       uint64
	isStarted          bool
	lastTalliedBlockNo uint64
}

func NewFavorChoice(threshold uint64) *FavorChoice {
	return &FavorChoice{
		Threshold:      threshold,
		voted:          make(map[uint64]uint64),
		records:        make(map[string]bool),
		selectedOption: types.NoOptionMade,
		isStarted:      false,
		noOfOption:     0,
	}
}

func (this FavorChoice) Desc() string {
	return `If a choice get more than ` + strconv.FormatUint(this.Threshold, 64) + ` then it got selected. 
Tally everytime there is a new vote and a person can only vote once.`
}

func (this *FavorChoice) Record(who string, raw interface{}) {
	option, _ := raw.(uint64)
	if this.records[who] == true {
		return
	}
	if this.selectedOption == types.NoOptionMade {
		this.records[who] = true
		this.voted[option] += 1
		this.lastVotedOpt = option
	}
}

func (this *FavorChoice) VotingPower(who string, option interface{}) uint64 {
	return 1
}

func (this *FavorChoice) Cost(who string, option interface{}) uint64 {
	return 0
}

func (this *FavorChoice) Tally(lastTalliedBlockNo uint64) bool {
	if this.selectedOption == types.NoOptionMade {
		if this.voted[this.lastVotedOpt] >= this.Threshold {
			this.selectedOption = this.lastVotedOpt
		}
	}
	this.lastTalliedBlockNo = lastTalliedBlockNo
	return true
}

func (this *FavorChoice) TallyAt() uint64 {
	return types.TallyAtVote
}

func (this *FavorChoice) GetTallyResult() ([]byte, uint64) {
	return nil, this.selectedOption
}

func (this *FavorChoice) Start(lastResult []byte, noOfOption uint64) bool {
	if this.isStarted == true {
		fmt.Println("FavorChoice already started!")
		return false
	}
	if lastResult != nil {
		buf := bytes.NewBuffer(lastResult)
		newThreshold, err := binary.ReadUint64(buf)
		if err == nil {
			this.Threshold = newThreshold
		}
	}
	if this.Threshold == 0 {
		return false
	} else {
		this.noOfOption = noOfOption
		this.isStarted = true
		return true
	}
}

func (this *FavorChoice) IsStarted() bool {
	return this.isStarted
}

func (this *FavorChoice) GetLastTalliedBlock() uint64 {
	return this.lastTalliedBlockNo
}

func (this *FavorChoice) ValidateVote(raw interface{}) bool {
	option, _ := raw.(uint64)
	if option < this.noOfOption {
		return true
	}
	return false
}
