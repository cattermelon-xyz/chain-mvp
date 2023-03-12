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

// TODO: change EventManager to interface instead of Struct to use in testing
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
	MissionId string
	Name      string
	Args      []byte
}

type Event struct {
	Name         string
	Args         []byte
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
				o.Update(e.Args)
			}
		}
		// log.Println("Emit ", id)
		this.Broadcast() <- EventData{
			Name: e.Name,
			Args: e.Args,
		}
		delete(this.registeredEvent, id)
	}
}

/**
* Params: name string, args []byte
* Return *Event, id string
 */
func (this *eventManagerStruct) CreateEvent(name string, args []byte) (*Event, string) {
	e := Event{Name: name, Args: args}
	id := utils.RandString(8)
	this.registeredEvent[id] = &e
	return &e, id
}
