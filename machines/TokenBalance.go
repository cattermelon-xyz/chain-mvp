package machines

import (
	"strconv"

	"github.com/hectagon-finance/chain-mvp/types"
)

// TODO: unfinished!
// In the 3 consecutive vote, if 2 choose a same Option then it will pass
// how to store data records[1]['abc']
type TokenBalance struct {
	Threshold    int
	records      map[int][]string
	power        map[string]int
	totalPower   int
	selected     int
	lastVotedOpt int
}

func NewTokenBalance(threshold int) *TokenBalance {
	return &TokenBalance{
		Threshold: threshold,
		records:   make(map[int][]string),
		selected:  types.NoTallyResult,
	}
}

func (this TokenBalance) Desc() string {
	return `If a choice get more than ` + strconv.Itoa(this.Threshold) + ` percent then it got selected. 
Tally everytime there is a new vote and a person can only as manytime as he like.`
}

func (this *TokenBalance) Record(who string, option int) {
	if this.selected == types.NoTallyResult {
		// &this.records[option] = append(&this.records[option],who)
		this.power[who] = this.VotingPower(who, option)
		this.lastVotedOpt = option
	}
}

func (this *TokenBalance) VotingPower(who string, option int) int {
	// query blockchain to get the power
	return 1
}

func (this *TokenBalance) Cost(who string, option int) int {
	return 0
}

func (this *TokenBalance) Tally() {
	if this.selected == types.NoTallyResult {
		// query blockchain to get totalPower
		// compare the last op
	}
}

func (this *TokenBalance) TallyAt() int {
	return types.TallyAtVote
}

func (this *TokenBalance) GetTallyResult() int {
	return this.selected
}
