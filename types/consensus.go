package types

import (
	"fmt"
	"time"
)

type BlockData struct {
	CurrentBlockNumber uint64
}

func GetCurrentBlockNumber() uint64 {
	BlockDataRequest <- true
	blockResp := <-BlockDataResp
	return blockResp.CurrentBlockNumber
}

var BlockDataResp = make(chan BlockData)
var BlockDataRequest = make(chan bool)
var CurrentBlockNumber uint64 = 1

func StartConsensus() {
	go produceBlock()
	go serveBlockRequest()
}

func produceBlock() {
	for {
		time.Sleep(time.Second * 3)
		CurrentBlockNumber += 1
		fmt.Println("CurrentBlockNo: ", CurrentBlockNumber)
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
