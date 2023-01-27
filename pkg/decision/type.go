package decision

import (
	"github.com/hectagon-finance/chain-mvp/pkg/checkpoint"
)

type Decision struct {
	Title    string
	Fulltext string
	Start    *checkpoint.CheckPoint
	Current  *checkpoint.CheckPoint
}
