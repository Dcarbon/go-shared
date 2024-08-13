package edef

import (
	"github.com/Dcarbon/go-shared/libs/ievent"
)

const (
	XNotfification = "notification"
)

type NotificationEvent struct {
	pusher ievent.IPublisher
}

func NewNotificationEvent(pusher ievent.IPublisher) *NotificationEvent {
	var ievent = &NotificationEvent{
		pusher: pusher,
	}
	return ievent
}

func (iotEvent *NotificationEvent) PushNotification(ev *EventPushNotification) error {
	return iotEvent.pusher.Publish(&ievent.Event{
		Exchange: XNotfification,
		Queue:    "",
		Data:     ev,
	})
}

type EventPushNotification struct {
	ProfileId string `json:"id"`
}
