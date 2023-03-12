package mock

type MockBlockchain struct {
	currentBlock uint64
}

func (this *MockBlockchain) GetCurrentBlockNumber() uint64 {
	return this.currentBlock
}

func (this *MockBlockchain) SetCurrentBlockNumber(n uint64) {
	this.currentBlock = n
}
