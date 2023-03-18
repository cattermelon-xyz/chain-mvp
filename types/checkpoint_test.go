package types_test

import (
	"fmt"
	"strconv"
	"testing"

	. "github.com/hectagon-finance/chain-mvp/mock"
	"github.com/hectagon-finance/chain-mvp/types"
	"github.com/stretchr/testify/assert"
)

type MockChkP struct {
	name       string
	children   []*types.CheckPoint
	FallbackId uint64
	want       types.CheckPointStartedStatus
}

func (m MockChkP) getChkP(mission *types.Mission) *types.CheckPoint {
	vm := MockVoteMachine{
		Started: true,
	}
	chkp := mission.CreateEmptyCheckPoint("test title", "test desc", &vm)
	chkp.FallbackId = m.FallbackId
	if len(m.children) > 0 {
		c := mission.CreateEmptyCheckPoint("child", "test", &vm)
		chkp.Attach(c.Id)
	}
	return chkp
}

func TestCheckPointStart(t *testing.T) {
	/**
	if len(this.children) == 0 || this.FallbackId == NoFallbackOption {
		return false
	}
	currentBlockNumber := this.blockchain.GetCurrentBlockNumber()
	return this.voteMachine.Start(lastTalliedResult, uint64(len(this.children)), currentBlockNumber, this.FallbackId)
	*/
	table := []MockChkP{
		{
			name:       fmt.Sprintf("No children & FallbackId is %s (const NoFallbackOption)", strconv.FormatUint(types.NoFallbackOption, 10)),
			children:   make([]*types.CheckPoint, 0),
			FallbackId: types.NoFallbackOption,
			want:       types.ChkPFailToStart,
		},
		{
			name: fmt.Sprintf("1 children & FallbackId is %s (const NoFallbackOption)", strconv.FormatUint(types.NoFallbackOption, 10)),
			children: []*types.CheckPoint{
				{},
			},
			FallbackId: types.NoFallbackOption,
			want:       types.ChkPFailToStart,
		},
		{
			name: fmt.Sprintf("1 children & FallbackId is %s", strconv.FormatUint(1, 10)),
			children: []*types.CheckPoint{
				{},
			},
			FallbackId: 1,
			want:       types.ChkPStarted,
		},
		{
			name:       fmt.Sprintf("No children & FallbackId is %s", strconv.FormatUint(types.NoFallbackOption, 10)),
			children:   make([]*types.CheckPoint, 0),
			FallbackId: 1,
			want:       types.ChkPFailToStart,
		},
	}
	m, _ := types.CreateMission("mission test", "desc", nil)
	for _, tt := range table {
		chkP := tt.getChkP(m)
		status, _ := types.ExportCheckPointStart(chkP, []byte{0})
		assert.Equal(t, status, tt.want, tt.name)
	}
}
