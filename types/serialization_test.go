package types_test

import (
	"testing"

	"github.com/hectagon-finance/chain-mvp/types"
)

func TestUnmarshal(t *testing.T) {
	md := &types.MissionData{
		Title:    "test",
		Fulltext: "test fulltext",
		Owner:    "test owner",
	}
	md.Unmarshal()
}
