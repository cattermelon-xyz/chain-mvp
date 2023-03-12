package event

import (
	"encoding/json"
	"log"
)

type EmitVotePredefinedEvent interface {
	EmitVoteRecorded(missionId string, who string)
	EmitVoteFailToRecord(missionId string, who string)
}

// VoteRecorded
func (this *eventManagerStruct) EmitVoteRecorded(missionId string, who string) {
	args := []string{missionId, who}
	b, _ := json.Marshal(args)
	log.Printf("%s - %s vote recorded\n", missionId, who)
	_, evid := this.CreateEvent(VoteRecorded, b)
	this.Emit(evid)
}

// VoteRecordFailed
func (this *eventManagerStruct) EmitVoteFailToRecord(missionId string, who string) {
	args := []string{missionId, who}
	b, _ := json.Marshal(args)
	log.Printf("%s - fail to record %s vote\n", missionId, who)
	_, evid := this.CreateEvent(VoteFailToRecord, b)
	this.Emit(evid)
}
