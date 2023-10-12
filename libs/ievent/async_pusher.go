package ievent

import (
	"context"

	"github.com/Dcarbon/go-shared/libs/rabbit"
)

// AsyncPusher :
type asyncPusher struct {
	*pusher
	evQueue chan *Event
	cancel  context.CancelFunc
}

// NewAsyncPusher :
func NewAsyncPusher(rbChan rabbit.IChannel, bandWidth int) IPublisher {
	var ap = &asyncPusher{
		pusher: &pusher{
			rbChan: rbChan,
		},
		evQueue: make(chan *Event, bandWidth),
	}

	ap.Start()
	return ap
}

// Start :
func (ap *asyncPusher) Start() {
	go ap.waitEvent()
}

// Stop :
func (ap *asyncPusher) Stop() {
	if nil != ap.cancel {
		ap.cancel()
		ap.cancel = nil
	}
}

func (ap *asyncPusher) Publish(ev *Event) error {
	ap.evQueue <- ev
	return nil
}

func (ap *asyncPusher) waitEvent() {
	if ap.cancel != nil {
		ap.cancel()
	}

	var ctx, cancel = context.WithCancel(context.Background())
	ap.cancel = cancel

	for {
		if nil == cancel {
			break
		}
		select {
		case ev := <-ap.evQueue:
			ap.push(ev)
		case <-ctx.Done():
			cancel()
			cancel = nil
		}
	}
}
