package event

import (
	"encoding/json"
	"log"
	"strconv"
)

type EmitCheckPointPredefinedEvent interface {
	EmitCheckPointStarted(missionId string, from string, to string)
	EmitCheckPointFailToStart(missionId string, fromChkPId string, selectedOption uint64)
}

// CheckPointStarted
func (this *eventManagerStruct) EmitCheckPointStarted(missionId string, fromChkPId string, toChkPId string) {
	args := []string{missionId, fromChkPId, toChkPId}
	b, _ := json.Marshal(args)
	log.Printf("%s - %s started (from %s)\n", missionId, toChkPId, fromChkPId)
	_, evid := this.CreateEvent(CheckPointStarted, b)
	this.Emit(evid)
}

// CheckPointFailToStart
func (this *eventManagerStruct) EmitCheckPointFailToStart(missionId string, fromChkPId string, selectedOption uint64) {
	args := []string{missionId, fromChkPId, strconv.Itoa(int(selectedOption))}
	b, _ := json.Marshal(args)
	log.Printf("%s - %s fail to choose option #%d\n", missionId, fromChkPId, selectedOption)
	_, evid := this.CreateEvent(CheckPointFailToStart, b)
	this.Emit(evid)
}
