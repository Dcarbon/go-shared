package edef

import (
	"github.com/Dcarbon/go-shared/dmodels"
	"github.com/Dcarbon/go-shared/libs/broker"
)

const (
	XIOTCreated      = "iot_created"
	XIOTChangeStatus = "iot_changed_status"
)

type IOTEvent struct {
	pusher broker.IPublisher
}

func NewIOTEvent(pusher broker.IPublisher) *IOTEvent {
	var ievent = &IOTEvent{
		pusher: pusher,
	}
	return ievent
}

func (ievent *IOTEvent) PushIOTCreate(ev *EventIOTCreate) error {
	return ievent.pusher.Publish(broker.Event{
		Exchange: XIOTCreated,
		Queue:    "",
		Data:     ev,
	})
}

func (ievent *IOTEvent) PushIOTChangeStatus(ev *EventIOTChangeStatus) error {
	return ievent.pusher.Publish(broker.Event{
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
