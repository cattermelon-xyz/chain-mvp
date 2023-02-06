package machines

import "github.com/hectagon-finance/chain-mvp/types"

//NOTE: this is not only rule but a VotingMachine with states stored in struct and functions

// if an option get more than 3 vote, it will pass
type ThreeVoteRuleData struct {
	Name  string
	Voted map[int]int
}

func (this ThreeVoteRuleData) Vote(tree *types.Initiative, who string, option int) {
	this.Voted[option] += 1
	if this.Voted[option] >= 3 {
		tree.Choose(option)
	}
}

func (this ThreeVoteRuleData) GetName() string {
	return this.Name
}

// in the 3 consecutive vote, if 2 choose a same option then it will pass
type FirstConsecutiveVoteRuleData struct {
	Name  string
	Voted []int
}

func (this FirstConsecutiveVoteRuleData) GetName() string {
	return this.Name
}

func (this FirstConsecutiveVoteRuleData) Vote(tree *types.Initiative, who string, option int) {

	if option == this.Voted[0] {
		// fmt.Printf("option : %d, Voted[0]: %d \n", option, this.Voted[0])
		tree.Choose(option)
	} else {
		this.Voted[2] = this.Voted[1]
		this.Voted[1] = this.Voted[0]
		this.Voted[0] = option
	}
}
