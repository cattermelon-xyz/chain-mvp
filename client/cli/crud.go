package main

import (
	"github.com/hectagon-finance/chain-mvp/types"
)

// create a decision, return empty if nothing is created
func create(who types.Address, title string, fulltext string) string

// do not override nil content
// TODO: use FLAG later
// func update(who net.Address, decisionId net.Address, newContent decision.Decision) bool

// delete Decision, what is the condition for deleting?
func delete(who types.Address, decisionId types.Address) bool
