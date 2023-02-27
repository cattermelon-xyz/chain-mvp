package types

import (
	"errors"

	"github.com/hectagon-finance/chain-mvp/third_party/utils"
)

type EventData struct {
	Name string
	Args []byte
}

var Broadcast = make(chan EventData) // channel to broadcast data

var registeredEvent = make(map[string]*Event)

type Observer interface {
	Update([]byte)
	GetId() string
	SetId(string)
}

type Event struct {
	Name         string
	Args         []byte
	observerList map[string]Observer
}

/**
* Emit(id string)
*	Emit event and notify all its Observer and Clear the Event from memory
* Params: id string
 */
func Emit(id string) {
	e := registeredEvent[id]
	if e != nil {
		if e.observerList != nil {
			for _, o := range e.observerList {
				o.Update(e.Args)
			}
		}
		// fmt.Println("Emit ", id)
		Broadcast <- EventData{
			Name: e.Name,
			Args: e.Args,
		}
		delete(registeredEvent, id)
	}
}

func CreateEvent(name string, args []byte) (*Event, string) {
	e := Event{Name: name, Args: args}
	id := utils.RandString(8)
	registeredEvent[id] = &e
	return &e, id
}

func DeleteEvent(eventId string) (bool, error) {
	if _, ok := registeredEvent[eventId]; ok {
		delete(registeredEvent, eventId)
		return true, nil
	}
	return false, errors.New(eventId + " not found")
}

func Register(eventId string, o Observer) (string, error) {
	e := registeredEvent[eventId]
	var oId string
	if e != nil {
		oId = utils.RandString(8)
		e.observerList[oId] = o
		o.SetId(oId)
		return oId, nil
	}
	return oId, errors.New(eventId + " not found")
}

func Deregister(eventId string, oId string) (bool, error) {
	e := registeredEvent[eventId]
	if e != nil {
		if _, ok := e.observerList[oId]; ok {
			delete(e.observerList, oId)
			return true, nil
		} else {
			return false, errors.New(oId + " not found")
		}
	}
	return false, errors.New(eventId + " not found")
}
