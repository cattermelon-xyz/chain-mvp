package net

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
