package types

import "log"

type Who struct {
	Network  Network
	Identity Address
}

func (this *Who) vote(tree *Mission, option []byte) {
	isValid := tree.IsValidChoice(option)
	if !isValid {
		log.Printf("%s vote %d, this is an invalid vote\n", this.Identity, option)
		return
	}
}
