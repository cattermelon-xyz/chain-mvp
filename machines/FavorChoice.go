package machines

import (
	"strconv"

	"github.com/hectagon-finance/chain-mvp/types"
)

type FavorChoice struct {
	Threshold    int
	voted        map[int]int
	records      map[string]bool
	selected     int
	lastVotedOpt int
}

func NewFavorChoice(threshold int) *FavorChoice {
	return &FavorChoice{
		Threshold: threshold,
		voted:     make(map[int]int),
		records:   make(map[string]bool),
		selected:  types.NoTallyResult,
	}
}

func (this FavorChoice) Desc() string {
	return `If a choice get more than ` + strconv.Itoa(this.Threshold) + ` then it got selected. 
Tally everytime there is a new vote and a person can only vote once.`
}

func (this *FavorChoice) Record(who string, option int) {
	if this.records[who] == true {
		return
	}
	if this.selected == types.NoTallyResult {
		this.records[who] = true
		this.voted[option] += 1
		this.lastVotedOpt = option
	}
}

func (this *FavorChoice) VotingPower(who string, option int) int {
	return 1
}

func (this *FavorChoice) Cost(who string, option int) int {
	return 0
}

func (this *FavorChoice) Tally() {
	if this.selected == types.NoTallyResult {
		if this.voted[this.lastVotedOpt] >= this.Threshold {
			this.selected = this.lastVotedOpt
		}
	}
}

func (this *FavorChoice) TallyAt() int {
	return types.TallyAtVote
}

func (this *FavorChoice) GetTallyResult() int {
	return this.selected
}
