package event

import (
	"encoding/json"
	"log"
)

type EmitFallbackPredefinedEvent interface {
	EmitFallbackAttempt(missionId string, checkPointId string)
}

func (this *eventManagerStruct) EmitFallbackAttempt(missionId string, checkPointId string) {
	args := []string{missionId, checkPointId}
	b, _ := json.Marshal(args)
	log.Printf("%s - CheckPoint %s attempt to fallback\n", missionId, checkPointId)
	_, evid := this.CreateEvent(FallbackAttempt, b)
	this.Emit(evid)
}
