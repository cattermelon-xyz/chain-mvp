package event

import (
	"encoding/json"
	"log"
)

type EmitTallyPredefinedEvent interface {
	EmitTallySucceed(missionId string, checkPointId string)
	EmitTallyFailed(missionId string, checkPointId string)
}

// TallySucceed
func (this *eventManagerStruct) EmitTallySucceed(missionId string, checkPointId string) {
	args := []string{missionId, checkPointId}
	b, _ := json.Marshal(args)
	log.Printf("%s - CheckPoint %s tallied\n", missionId, checkPointId)
	_, evid := this.CreateEvent(TallySucceed, b)
	this.Emit(evid)
}

// TallyFailed
func (this *eventManagerStruct) EmitTallyFailed(missionId string, checkPointId string) {
	args := []string{missionId, checkPointId}
	b, _ := json.Marshal(args)
	log.Printf("%s - CheckPoint %s fail to tally\n", missionId, checkPointId)
	_, evid := this.CreateEvent(TallyFailed, b)
	this.Emit(evid)
}
