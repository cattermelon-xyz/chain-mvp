package main

import (
	"fmt"

	"github.com/hectagon-finance/chain-mvp/machines"
	"github.com/hectagon-finance/chain-mvp/types"
)

func main() {
	dRule := machines.ThreeVoteRuleData{Voted: make(map[int]int), Name: "d_ThreeVoteRuleData"}
	cRule := machines.FirstConsecutiveVoteRuleData{Voted: []int{-1, -1, -1}, Name: "c_FirstConsecutiveVoteRuleData"}
	bRule := machines.ThreeVoteRuleData{Voted: make(map[int]int), Name: "b_ThreeVoteRuleData"}
	aRule := machines.FirstConsecutiveVoteRuleData{Voted: []int{-1, -1, -1}, Name: "a_FirstConsecutiveVoteRuleData"}

	e, f, g, h, k := types.CreateEmptyNode("e", nil), types.CreateEmptyNode("f", nil), types.CreateEmptyNode("g", nil), types.CreateEmptyNode("h", nil), types.CreateEmptyNode("k", nil)
	d := types.CreateNodeWithChildren("d", []*types.Node{g, h, k}, dRule)
	c := types.CreateNodeWithChildren("c", []*types.Node{f}, cRule)
	b := types.CreateNodeWithChildren("b", []*types.Node{d, e}, bRule)
	a := types.CreateNodeWithChildren("a", []*types.Node{b, c}, aRule)
	t, _ := types.CreateInitiative("DEMO", "A simple Initiative", a)
	fmt.Println("Original initiative:")
	t.PrintFromCurrent()
	t.Start()
	t.Vote(0, "alice")
	t.Vote(1, "bob")
	t.Pause()
	t.Vote(1, "bob1")
	t.Vote(1, "caite")
	t.Resume()
	t.Vote(1, "david")
	t.Vote(1, "erh")
	t.Vote(0, "felix")
	t.Vote(1, "garin")
}
