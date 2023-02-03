package datastrct

import (
	"fmt"

	"github.com/hectagon-finance/chain-mvp/third_party/tree"
)

type Tree struct {
	Start   *Node
	Current *Node
}

func CreateTree(root *Node) *Tree {
	tree := Tree{
		Start:   root,
		Current: root,
	}
	return &tree
}
func (this *Tree) print() {
	tree.Print(this.Start)
}
func (this *Tree) PrintFromCurrent() {
	tree.Print(this.Current)
}
func (this *Tree) Choose(idx int) {
	nextNode := this.Current.Get(idx)
	if nextNode == nil {
		fmt.Println(idx, " out of bound, no move")
	}
	if nextNode != nil {
		fmt.Printf("from %s choose: %d got %s\n", this.Current.name, idx, nextNode.name)
		this.Current = nextNode
		this.PrintFromCurrent()
	}
}
func (this *Tree) IsValidChoice(idx int) bool {
	return this.Current.isValidChoice(idx)
}
func (this *Tree) Vote(idx int, who string) {
	if this.IsValidChoice(idx) {
		fmt.Printf("In %s, with %s, %s vote %d\n", this.Current.Data(), this.Current.voteObject.name(), who, idx)
		this.Current.vote(this, who, idx)
	}
}
