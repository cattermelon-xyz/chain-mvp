package event

import (
	"encoding/json"
	"log"
)

type EmitMissionPredefinedEvent interface {
	EmitMissionStarted(missionId string)
	EmitMissionPaused(missionId string)
	EmitMissionResumed(missionId string)
	EmitMissionStopped(missionId string)
}

func (this eventManagerStruct) EmitMissionStarted(missionId string) {
	args := []string{missionId}
	b, _ := json.Marshal(args)
	log.Printf("%s - Mission started successfully\n", missionId)
	_, evid := this.CreateEvent(MissionStarted, b)
	this.Emit(evid)
}

func (this eventManagerStruct) EmitMissionPaused(missionId string) {
	args := []string{missionId}
	b, _ := json.Marshal(args)
	log.Printf("%s - Mission paused \n", missionId)
	_, evid := this.CreateEvent(MissionPaused, b)
	this.Emit(evid)
}

func (this eventManagerStruct) EmitMissionResumed(missionId string) {
	args := []string{missionId}
	b, _ := json.Marshal(args)
	log.Printf("%s - Mission resumed \n", missionId)
	_, evid := this.CreateEvent(MissionResumed, b)
	this.Emit(evid)
}

func (this eventManagerStruct) EmitMissionStopped(missionId string) {
	args := []string{missionId}
	b, _ := json.Marshal(args)
	log.Printf("%s - Mission stopped \n", missionId)
	_, evid := this.CreateEvent(MissionStopped, b)
	this.Emit(evid)
}
