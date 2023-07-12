package broker

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Dcarbon/go-shared/libs/rabbit"
	"github.com/streadway/amqp"
)

// AsyncPusher :
type asyncPusher struct {
	evQueue chan Event
	rbChan  rabbit.IChannel
	cancel  context.CancelFunc
}

// NewAsyncPusher :
func NewAsyncPusher(rbChan rabbit.IChannel, bandWidth int) IPublisher {
	var ap = &asyncPusher{
		evQueue: make(chan Event, bandWidth),
		rbChan:  rbChan,
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

func (ap *asyncPusher) Publish(ev Event) error {
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
			ap.push(&ev)
		case <-ctx.Done():
			cancel()
			cancel = nil
		}
	}
}

func (ap *asyncPusher) push(ev *Event) {
	var raw []byte
	var err error

	if tmp, ok := ev.Data.([]byte); ok {
		raw = tmp
	} else {
		raw, err = json.Marshal(ev.Data)
		if nil != err {
			log.Println("asyncPusher: marshall event error: ", ev, err)
			return
		}
	}

	err = ap.rbChan.Publish(
		ev.Exchange,
		ev.Queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        raw,
			Headers:     ev.Headers,
		})

	if nil != err {
		log.Printf("Push event to queue %s error: %s\n", ev.Queue, err.Error())
	}
}
