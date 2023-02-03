package main

import "github.com/hectagon-finance/chain-mvp/pkg/net"

// this binary will handle vote function
// user: organization users
// what params should this binary take?
/*
	// vote for a decision in its current CheckPoint
	// return true/false dictate if the vote valid and recorded
	vote(who *net.Who, decisionId string) bool
*/

func vote(whoAddr net.Address, decisionId net.Address)
