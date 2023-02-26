package types

type Network struct {
	Title    string
	Version  string
	Endpoint []string
}

type Address string

var dchain Network = Network{
	Title:    "DChain",
	Version:  "0.0t",
	Endpoint: []string{"endpoint1"},
}

var ethereum Network = Network{
	Title:    "Ethereum",
	Version:  "1p",
	Endpoint: []string{"endpoint1"},
}

var solana Network = Network{
	Title:    "Solana",
	Version:  "1p",
	Endpoint: []string{"endpoint1"},
}

func StringToAddress(str string) Address {
	// check
	return Address(str)
}

func GetNetWork(networkId string) Network {
	var result Network
	switch networkId {
	case "DChain":
		result = dchain
		break
	case "ETH":
		result = ethereum
		break
	case "SOL":
		result = solana
		break
	}
	return result
}

func GetCurrentBlockNumber() uint64 {
	return 100
}

// func CreateAddress() Address {
// 	return StringToAddress(utils.randStringBytesMaskImpr(16))
// }
