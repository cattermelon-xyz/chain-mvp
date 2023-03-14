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

func (m MockChkP) getChkP() *types.CheckPoint {
	chkp := types.CreateEmptyCheckPoint("test title", "test desc", &MockVoteMachine{
		Started: true,
	}, &MockBlockchain{})

	if len(m.children) > 0 {
		chkp.Attach(m.children[0])
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

	for _, tt := range table {
		chkP := tt.getChkP()

		assert.Equal(t, types.ExportCheckPointStart(chkP, []byte{0}), tt.want, tt.name)
	}
}
