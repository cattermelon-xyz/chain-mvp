package main

import (
	"log"

	"github.com/hectagon-finance/chain-mvp/machines"
	"github.com/hectagon-finance/chain-mvp/types"
)

func main() {
	dRule := machines.NewFavorChoice(3)
	cRule := machines.NewFavorChoice(3)
	bRule := machines.NewFavorChoice(3)
	aRule := machines.NewFavorChoice(3)

	e, f, g, h, k := types.CreateEmptyNode("e", nil), types.CreateEmptyNode("f", nil), types.CreateEmptyNode("g", nil), types.CreateEmptyNode("h", nil), types.CreateEmptyNode("k", nil)
	d := types.CreateNodeWithChildren("d", []*types.Node{g, h, k}, dRule)
	c := types.CreateNodeWithChildren("c", []*types.Node{f}, cRule)
	b := types.CreateNodeWithChildren("b", []*types.Node{d, e}, bRule)
	a := types.CreateNodeWithChildren("a", []*types.Node{b, c}, aRule)
	t, _ := types.CreateInitiative("DEMO", "A simple Initiative", a)
	log.Println("Original initiative:")
	t.PrintFromCurrent()
	t.Start()
	t.Vote(0, "alice")
	t.Vote(0, "bob")
	// t.Pause()
	t.Vote(1, "bob1")
	t.Vote(0, "caite")
	// t.Resume()
	t.Vote(1, "david")
	t.Vote(1, "erh")
	t.Vote(0, "felix")
	t.Vote(1, "garin")
}
