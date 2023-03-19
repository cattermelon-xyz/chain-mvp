package event

import (
	"errors"

	"github.com/hectagon-finance/chain-mvp/third_party/utils"
)

// One Observer can Register to only ONE event (through SetId)
type Observer interface {
	Update([]byte)
	GetId() string
	SetId(string)
	Marshal() []byte
	Unmarshal(oId string, data []byte)
}

func (this *eventManagerStruct) DeleteEvent(eventId string) (bool, error) {
	if _, ok := this.registeredEvent[eventId]; ok {
		delete(this.registeredEvent, eventId)
		return true, nil
	}
	return false, errors.New(eventId + " not found")
}

func (this *eventManagerStruct) Register(eventId string, o Observer) (string, error) {
	e := this.registeredEvent[eventId]
	var oId string
	if e != nil {
		oId = utils.RandString(8)
		e.observerList[oId] = o
		o.SetId(oId)
		return oId, nil
	}
	return oId, errors.New(eventId + " not found")
}

func (this *eventManagerStruct) Deregister(eventId string, oId string) (bool, error) {
	e := this.registeredEvent[eventId]
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
