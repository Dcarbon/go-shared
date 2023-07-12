package broker

import (
	"encoding/json"
	"log"

	"github.com/Dcarbon/go-shared/libs/rabbit"
	"github.com/streadway/amqp"
)

// directPusher :
type directPusher struct {
	rbChan rabbit.IChannel
}

// NewAsyncPusher :
func NewDirectPusher(rbChan rabbit.IChannel) IPublisher {
	var ap = &directPusher{
		rbChan: rbChan,
	}
	return ap
}

func (ap *directPusher) Publish(ev Event) error {
	return ap.push(&ev)
}

func (ap *directPusher) push(ev *Event) error {
	var raw []byte
	var err error

	if tmp, ok := ev.Data.([]byte); ok {
		raw = tmp
	} else {
		raw, err = json.Marshal(ev.Data)
		if nil != err {
			log.Println("directPusher: marshall event error: ", ev, err)
			return err
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
	return err
}
