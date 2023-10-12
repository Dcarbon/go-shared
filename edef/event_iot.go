package edef

import (
	"github.com/Dcarbon/go-shared/dmodels"
	"github.com/Dcarbon/go-shared/libs/ievent"
)

const (
	XIOTCreated      = "iot_created"
	XIOTChangeStatus = "iot_changed_status"
)

type IOTEvent struct {
	pusher ievent.IPublisher
}

func NewIOTEvent(pusher ievent.IPublisher) *IOTEvent {
	var ievent = &IOTEvent{
		pusher: pusher,
	}
	return ievent
}

func (iotEvent *IOTEvent) PushIOTCreate(ev *EventIOTCreate) error {
	return iotEvent.pusher.Publish(&ievent.Event{
		Exchange: XIOTCreated,
		Queue:    "",
		Data:     ev,
	})
}

func (iotEvent *IOTEvent) PushIOTChangeStatus(ev *EventIOTChangeStatus) error {
	return iotEvent.pusher.Publish(&ievent.Event{
		Exchange: XIOTChangeStatus,
		Queue:    "",
		Data:     ev,
	})
}

type GPS struct {
	Lat float64 `json:"lat"` // vi tuyen (pgis: y)
	Lng float64 `json:"lng"` // kinh tuyen:(pgis: x)
}

type EventIOTCreate struct {
	ID       int64                `json:"id"`
	Status   dmodels.DeviceStatus `json:"status"`
	Address  string               `json:"address"`
	Location *GPS                 `json:"location"`
}

type EventIOTChangeStatus struct {
	ID     int64                `json:"id"`
	Status dmodels.DeviceStatus `json:"status"`
}
