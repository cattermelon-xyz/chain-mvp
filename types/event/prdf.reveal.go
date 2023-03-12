package event

import (
	"encoding/json"
	"log"
)

type EmitRevealEvent interface {
	EmitRevealSucceed(missionId string, who string)
}

// TODO: should emit private key?
func (this *eventManagerStruct) EmitRevealSucceed(missionId string, who string) {
	args := []string{missionId, who}
	b, _ := json.Marshal(args)
	log.Printf("%s - %s revealed the vote result\n", missionId, who)
	_, evid := this.CreateEvent(RevealSucceed, b)
	this.Emit(evid)
}
