package decision

import (
	"github.com/hectagon-finance/chain-mvp/pkg/checkpoint"
	"github.com/hectagon-finance/chain-mvp/pkg/net"
)

type Decision struct {
	Title    string
	Fulltext string
	Start    *checkpoint.CheckPoint
	Id       net.Address
	Owner    net.Address
	State    *State
}

type State struct {
	IsStarted bool
	Current   *checkpoint.CheckPoint
}

type iDecision interface {
	create(Decision) string
}

func create(title string, fulltext string, start *checkpoint.CheckPoint) *Decision

func (this *Decision) edit(d Decision) bool

func (this *Decision) delete(d Decision) bool

func (this *Decision) start(d Decision) bool

func (this *Decision) stop(d Decision) bool

func (this *Decision) pause(d Decision) bool

func (this *Decision) resume(d Decision) bool
