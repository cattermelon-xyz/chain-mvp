package main

import (
	"fmt"

	"github.com/hectagon-finance/chain-mvp/pkg/datastrct"
)

func main() {
	dRule := datastrct.ThreeVoteRuleData{Voted: make(map[int]int), Name: "d_ThreeVoteRuleData"}
	cRule := datastrct.FirstConsecutiveVoteRuleData{Voted: []int{-1, -1, -1}, Name: "c_FirstConsecutiveVoteRuleData"}
	bRule := datastrct.ThreeVoteRuleData{Voted: make(map[int]int), Name: "b_ThreeVoteRuleData"}
	aRule := datastrct.FirstConsecutiveVoteRuleData{Voted: []int{-1, -1, -1}, Name: "a_FirstConsecutiveVoteRuleData"}

	e, f, g, h, k := datastrct.CreateEmptyNode("e", nil), datastrct.CreateEmptyNode("f", nil), datastrct.CreateEmptyNode("g", nil), datastrct.CreateEmptyNode("h", nil), datastrct.CreateEmptyNode("k", nil)
	d := datastrct.CreateNodeWithChildren("d", []*datastrct.Node{g, h, k}, dRule)
	c := datastrct.CreateNodeWithChildren("c", []*datastrct.Node{f}, cRule)
	b := datastrct.CreateNodeWithChildren("b", []*datastrct.Node{d, e}, bRule)
	a := datastrct.CreateNodeWithChildren("a", []*datastrct.Node{b, c}, aRule)
	t := datastrct.CreateTree(a)
	fmt.Println("Original tree:")
	t.PrintFromCurrent()
	t.Vote(0, "alice")
	t.Vote(0, "bob")
	t.Vote(1, "caite")
	t.Vote(1, "david")
}
