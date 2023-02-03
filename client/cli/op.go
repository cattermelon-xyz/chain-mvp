package main

import (
	"github.com/hectagon-finance/chain-mvp/types"
)

func getDChainWho(addr types.Address) types.Who {
	return types.Who{
		Network:  types.GetNetWork("DChain"),
		Identity: addr,
	}
}

func isDecisioExisted(decisionId string) bool

// print the detail of the Decision tree with current CheckPoint highlighted
func print() {}

// start, return true/false if decision is successfully started
func start(whoAddr types.Address, decisionId types.Address) bool {
	// TODO: convert whoAddr to who of this network -> guess this is what cosmos provide
	// who := getDChainWho(whoAddr)
	// TODO: check if decision is existed
	return true
}

// stop, return true/false if decision is successfully stopped
func stop(who types.Address, decisionId types.Address) bool

// pause, return true/false if decision is successfully paused
func pause(who types.Address, decisionId types.Address) bool

// resume, return true/false if decision is successfully started
func resume(who types.Address, decisionId types.Address) bool
