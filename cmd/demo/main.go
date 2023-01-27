package main

import (
	"fmt"

	"github.com/hectagon-finance/chain-mvp/pkg/checkpoint"
	"github.com/hectagon-finance/chain-mvp/pkg/decision"
	"github.com/hectagon-finance/chain-mvp/pkg/enforcer"
	"github.com/hectagon-finance/chain-mvp/pkg/net"
)

func main() {
	ETH := net.Network{Title: "ethereum", Version: "london", Endpoint: []string{"e_endpoint1", "e_endpoint2"}}
	// SOL := Network{"solana", "sealevel", []string{"s_endpoint1", "s_endpoint2"}}
	disperseContractAddress := "0xDisperseContract"
	params := []string{"0xFrom", "0xTo"}

	// checkpoint1 -> checkpoint2 -> enforce

	vote2ContractAddress := "0xVote2"
	checkpoint2VoteFunc := checkpoint.VoteFunc{Network: ETH, Address: vote2ContractAddress}
	enforceSendEthToAddress := enforcer.Enforcer{Network: ETH, Address: disperseContractAddress, Params: params}
	output2 := checkpoint.VoteOutput{Next: nil, Enforcer: &enforceSendEthToAddress}
	output2s := []checkpoint.VoteOutput{output2}

	checkpoint2 := checkpoint.CheckPoint{
		Title: "CheckPoint 2", Fulltext: "CheckPoint 2 FullText",
		Options: &output2s, VoteFunc: &checkpoint2VoteFunc}
	vote1ContractAddress := "0xVote1"
	checkpoint1VoteFunc := checkpoint.VoteFunc{Network: ETH, Address: vote1ContractAddress}
	output1 := checkpoint.VoteOutput{Next: &checkpoint2, Enforcer: nil}
	output1s := []checkpoint.VoteOutput{output1}
	checkpoint1 := checkpoint.CheckPoint{
		Title: "CheckPoint 1", Fulltext: "CheckPoint 1 FullText",
		Options: &output1s, VoteFunc: &checkpoint1VoteFunc}

	d := decision.Decision{Title: "Decision1", Fulltext: "Fulltext Decision 1",
		Start: &checkpoint1, Current: nil}

	fmt.Println("Hello World ", d.Title)
}
