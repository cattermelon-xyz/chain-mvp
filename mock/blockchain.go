package mock

import "github.com/hectagon-finance/chain-mvp/types/event"

type MockBlockchain struct {
	currentBlock     uint64
	mockEventManager *MockEventManager
}

func (this *MockBlockchain) GetCurrentBlockNumber() uint64 {
	return this.currentBlock
}

func (this *MockBlockchain) SetCurrentBlockNumber(n uint64) {
	this.currentBlock = n
}

func (this *MockBlockchain) GetEventManager() event.EventManager {
	if this.mockEventManager == nil {
		this.mockEventManager = &MockEventManager{}
	}
	return this.mockEventManager
}

func (this *MockBlockchain) SetEventManager(m *MockEventManager) {
	this.mockEventManager = m
}
