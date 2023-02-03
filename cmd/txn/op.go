package main

import (
	"github.com/hectagon-finance/chain-mvp/pkg/net"
)

func getDChainWho(addr net.Address) net.Who {
	return net.Who{
		Network:  net.GetNetWork("DChain"),
		Identity: addr,
	}
}

func isDecisioExisted(decisionId string) bool

// print the detail of the Decision tree with current CheckPoint highlighted
func print() {}

// start, return true/false if decision is successfully started
func start(whoAddr net.Address, decisionId net.Address) bool {
	// TODO: convert whoAddr to who of this network -> guess this is what cosmos provide
	// who := getDChainWho(whoAddr)
	// TODO: check if decision is existed
	return true
}

// stop, return true/false if decision is successfully stopped
func stop(who net.Address, decisionId net.Address) bool

// pause, return true/false if decision is successfully paused
func pause(who net.Address, decisionId net.Address) bool

// resume, return true/false if decision is successfully started
func resume(who net.Address, decisionId net.Address) bool
