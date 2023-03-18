package mock

import (
	"github.com/hectagon-finance/chain-mvp/types/event"
)

type MockEventManager struct {
	queue []string
}

func (this *MockEventManager) GetQueue() []string {
	return this.queue
}

func (this *MockEventManager) Queue(ev string) {
	if this.queue == nil {
		this.queue = make([]string, 0)
	}
	this.queue = append(this.queue, ev)
}

func (this *MockEventManager) UnQueue() string {
	if this.queue == nil || len(this.queue) == 0 {
		return ""
	}
	rs := this.queue[0]
	this.queue = this.queue[1:]
	return rs
}

func (this MockEventManager) Clear() {
	this.queue = make([]string, 0)
	return
}
func (this *MockEventManager) CreateEvent(name string, args []byte) (*event.Event, string) {
	return &event.Event{
		Id:   name,
		Name: name,
		Args: args,
	}, name
}
func (this *MockEventManager) DeleteEvent(eventId string) (bool, error) {
	return false, nil
}
func (this *MockEventManager) Register(eventId string, o event.Observer) (string, error) {
	return "", nil
}
func (this *MockEventManager) Deregister(eventId string, oId string) (bool, error) {
	return false, nil
}
func (this *MockEventManager) Emit(id string) {
	this.Queue(id)
}
func (this *MockEventManager) Broadcast() chan event.Event {
	return nil
}
func (this *MockEventManager) EmitMissionStarted(missionId string) {
	this.Queue(event.MissionStarted)
}
func (this *MockEventManager) EmitMissionPaused(missionId string) {
	this.Queue(event.MissionPaused)
}
func (this *MockEventManager) EmitMissionResumed(missionId string) {
	this.Queue(event.MissionResumed)
}
func (this *MockEventManager) EmitMissionStopped(missionId string) {
	this.Queue(event.MissionStopped)
}
func (this *MockEventManager) EmitCheckPointStarted(missionId string, fromChkPId string, toChkPId string) {
	this.Queue(event.CheckPointStarted)
}
func (this *MockEventManager) EmitCheckPointFailToStart(missionId string, fromChkPId string, selectedOption uint64) {
	this.Queue(event.CheckPointFailToStart)
}
func (this *MockEventManager) EmitFallbackAttempt(missionId string, checkPointId string) {
	this.Queue(event.FallbackAttempt)
}
func (this *MockEventManager) EmitRevealSucceed(missionId string, who string) {
	this.Queue(event.RevealSucceed)
}
func (this *MockEventManager) EmitTallySucceed(missionId string, checkPointId string) {
	this.Queue(event.TallySucceed)
}
func (this *MockEventManager) EmitTallyFailed(missionId string, checkPointId string) {
	this.Queue(event.TallyFailed)
}
func (this *MockEventManager) EmitVoteRecorded(missionId string, who string) {
	this.Queue(event.VoteRecorded)
}
func (this *MockEventManager) EmitVoteFailToRecord(missionId string, who string) {
	this.Queue(event.VoteFailToRecord)
}
