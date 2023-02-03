package main

import (
	"fmt"

	"github.com/hectagon-finance/chain-mvp/rules"
	"github.com/hectagon-finance/chain-mvp/types"
)

func main() {
	dRule := rules.ThreeVoteRuleData{Voted: make(map[int]int), Name: "d_ThreeVoteRuleData"}
	cRule := rules.FirstConsecutiveVoteRuleData{Voted: []int{-1, -1, -1}, Name: "c_FirstConsecutiveVoteRuleData"}
	bRule := rules.ThreeVoteRuleData{Voted: make(map[int]int), Name: "b_ThreeVoteRuleData"}
	aRule := rules.FirstConsecutiveVoteRuleData{Voted: []int{-1, -1, -1}, Name: "a_FirstConsecutiveVoteRuleData"}

	e, f, g, h, k := types.CreateEmptyNode("e", nil), types.CreateEmptyNode("f", nil), types.CreateEmptyNode("g", nil), types.CreateEmptyNode("h", nil), types.CreateEmptyNode("k", nil)
	d := types.CreateNodeWithChildren("d", []*types.Node{g, h, k}, dRule)
	c := types.CreateNodeWithChildren("c", []*types.Node{f}, cRule)
	b := types.CreateNodeWithChildren("b", []*types.Node{d, e}, bRule)
	a := types.CreateNodeWithChildren("a", []*types.Node{b, c}, aRule)
	t := types.CreateTree(a)
	fmt.Println("Original tree:")
	t.PrintFromCurrent()
	t.Vote(0, "alice")
	t.Vote(0, "bob")
	t.Vote(1, "caite")
	t.Vote(1, "david")
}
