package types

import (
	"log"
	"time"

	"github.com/hectagon-finance/chain-mvp/types/event"
)

type blockchainStruct struct {
	isInitialized bool
}

type Blockchain interface {
	GetCurrentBlockNumber() uint64
	GetEventManager() event.EventManager
}

type BlockData struct {
	CurrentBlockNumber uint64
}

var globalBlockchain = blockchainStruct{
	isInitialized: false,
}

/**
* Singleton method to return the Blockchain
 */
func GetBlockchain() Blockchain {
	if globalBlockchain.isInitialized == false {
		globalBlockchain = blockchainStruct{
			isInitialized: true,
		}
	}
	return &globalBlockchain
}

func (this *blockchainStruct) GetCurrentBlockNumber() uint64 {
	BlockDataRequest <- true
	blockResp := <-BlockDataResp
	return blockResp.CurrentBlockNumber
}

func (this *blockchainStruct) GetEventManager() event.EventManager {
	return event.GetEventManager()
}

var BlockDataResp = make(chan BlockData)
var BlockDataRequest = make(chan bool)
var CurrentBlockNumber uint64 = 1

func StartConsensus(m event.EmitPredefinedEvent) {
	go produceBlock(m)
	go serveBlockRequest()
}

func produceBlock(m event.EmitPredefinedEvent) {
	for {
		time.Sleep(time.Second * 3)
		CurrentBlockNumber += 1
		startBlock(m)
		log.Println("CurrentBlockNo: ", CurrentBlockNumber)
	}
}

/**
* Func startBlock
* Loop through all missions and perform tally
 */
func startBlock(ev event.EmitPredefinedEvent) {
	for _, m := range Missions {
		if m.isActive == true {
			m.BeatAtNewBlock()
		}
	}
}

func serveBlockRequest() {
	for {
		request := <-BlockDataRequest
		if request == true {
			BlockDataResp <- BlockData{
				CurrentBlockNumber: CurrentBlockNumber,
			}
		}
	}
}
