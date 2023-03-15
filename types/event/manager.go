package event

import "github.com/hectagon-finance/chain-mvp/third_party/utils"

type PredefinedEventName string

const (
	MissionStarted        = "MissionStarted"
	MissionPaused         = "MissionPaused"
	MissionResumed        = "MissionResumed"
	MissionStopped        = "MissionStopped"
	CheckPointStarted     = "CheckPointStarted"
	CheckPointFailToStart = "CheckPointFailToStart"
	FallbackAttempt       = "FallbackAttempt"
	RevealSucceed         = "RevealSucceed"
	TallySucceed          = "TallySucceed"
	TallyFailed           = "TallyFailed"
	VoteRecorded          = "VoteRecorded"
	VoteFailToRecord      = "VoteFailToRecord"
)

type eventManagerStruct struct {
	initialized     bool
	registeredEvent map[string]*Event
	broadcast       chan EventData
}

var globalEventManager = eventManagerStruct{
	initialized: false,
}

type EmitPredefinedEvent interface {
	EmitCheckPointPredefinedEvent
	EmitFallbackPredefinedEvent
	EmitMissionPredefinedEvent
	EmitTallyPredefinedEvent
	EmitVotePredefinedEvent
	EmitRevealEvent
}

type EventManager interface {
	EmitPredefinedEvent
	Emit(string)
	CreateEvent(string, []byte) (*Event, string)
	DeleteEvent(string) (bool, error)
	Register(string, Observer) (string, error)
	Deregister(string, string) (bool, error)
	Broadcast() chan EventData
}

type EventData struct {
	Name string
	Args []byte
}

type Event struct {
	Id           string
	Data         EventData
	observerList map[string]Observer
}

/**
* Singleton, return global object
 */
func GetEventManager() EventManager {
	if globalEventManager.initialized == false {
		globalEventManager = eventManagerStruct{
			initialized:     true,
			registeredEvent: make(map[string]*Event),
			broadcast:       make(chan EventData),
		}
	}
	return EventManager(&globalEventManager)
}

func (this *eventManagerStruct) Broadcast() chan EventData {
	return this.broadcast
}

/**
* Emit(id string),
*	Emit event and notify all its Observer and Clear the Event from memory,
* Params: id string
 */
func (this *eventManagerStruct) Emit(id string) {
	e := this.registeredEvent[id]
	if e != nil {
		if e.observerList != nil {
			for _, o := range e.observerList {
				o.Update(e.Data.Args)
			}
		}
		// log.Println("Emit ", id)
		this.Broadcast() <- e.Data
		delete(this.registeredEvent, id)
	}
}

/**
* Params: name string, args []byte
* Return *Event, id string
 */
func (this *eventManagerStruct) CreateEvent(name string, args []byte) (*Event, string) {
	id := utils.RandString(8)
	e := Event{
		Data: EventData{Name: name, Args: args}, Id: id,
	}
	this.registeredEvent[id] = &e
	return &e, id
}
