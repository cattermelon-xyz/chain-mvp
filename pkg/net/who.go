// EOA - externally owned account

package net

import (
	"fmt"

	"github.com/hectagon-finance/chain-mvp/pkg/datastrct"
)

type Who struct {
	Network  Network
	Identity Address
}

func (this *Who) vote(tree *datastrct.Tree, option int) {
	isValid := tree.IsValidChoice(option)
	if !isValid {
		fmt.Printf("%s vote %d, this is an invalid vote\n", this.Identity, option)
		return
	}
}
