package main

import (
	"os"

	"github.com/hectagon-finance/chain-mvp/pkg/net"
)

// this binary will take create / TODO:update / delete command toward Decision
// user: organization operators
// what params should this binary take?
func main() {
	actor := net.StringToAddress(os.Args[0])
	cmd := os.Args[1]
	switch cmd {
	case "create":
		title := os.Args[2]
		fulltext := os.Args[3]
		create(actor, title, fulltext)
		break
	// case "update":
	// 	decisionId := net.StringToAddress(os.Args[2])
	// 	title := os.Args[3]
	// 	fulltext := os.Args[4]
	// 	update(actor, decisionId, decision.Decision{Title: title, Fulltext: fulltext,
	// 		Start:   nil, // TODO: placeholder
	// 		Current: nil})
	// 	break
	case "delete":
		decisionId := net.StringToAddress(os.Args[2])
		delete(actor, decisionId)
		break
	case "start":
		decisionId := net.StringToAddress(os.Args[2])
		start(actor, decisionId)
		break
	case "stop":
		decisionId := net.StringToAddress(os.Args[2])
		stop(actor, decisionId)
		break
	case "pause":
		decisionId := net.StringToAddress(os.Args[2])
		pause(actor, decisionId)
		break
	case "resume":
		decisionId := net.StringToAddress(os.Args[2])
		resume(actor, decisionId)
		break
	case "vote":
		decisionId := net.StringToAddress(os.Args[2])
		vote(actor, decisionId)
		break
	}
}
