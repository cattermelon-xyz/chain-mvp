package types

import (
	"fmt"
)

type Who struct {
	Network  Network
	Identity Address
}

func (this *Who) vote(tree *Mission, option int) {
	isValid := tree.IsValidChoice(option)
	if !isValid {
		fmt.Printf("%s vote %d, this is an invalid vote\n", this.Identity, option)
		return
	}
}
