package enforcer

import "github.com/hectagon-finance/chain-mvp/pkg/net"

type Enforcer struct {
	Network net.Network
	Address string
	Params  []string
}
