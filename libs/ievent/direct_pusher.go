package ievent

import (
	"log"

	"github.com/Dcarbon/go-shared/libs/rabbit"
	"github.com/streadway/amqp"
)

type pusher struct {
	rbChan rabbit.IChannel
}

func (p *pusher) push(ev *Event) error {
	var raw, err = ev.Encode()
	if nil != err {
		return err
	}

	var msg = amqp.Publishing{
		ContentType: "application/json",
		Body:        raw,
		Headers:     ev.Headers,
	}

	// if ev.Delay > 0 {
	// 	msg.Expiration = fmt.Sprintf("%d", ev.Delay)
	// }

	err = p.rbChan.Publish(
		ev.Exchange,
		ev.Queue,
		false,
		false,
		msg)
	if nil != err {
		log.Printf("Push event to queue %s error: %s\n",
			ev.Queue, err.Error(),
		)
	}

	return err
}

func (p *pusher) GetChannel() rabbit.IChannel {
	return p.rbChan
}

// directPusher :
type directPusher struct {
	*pusher
}

// NewAsyncPusher :
func NewDirectPusher(rbChan rabbit.IChannel,
) IPublisher {
	var ap = &directPusher{
		pusher: &pusher{
			rbChan: rbChan,
		},
	}
	return ap
}

func (ap *directPusher) Publish(ev *Event) error {
	return ap.push(ev)
}
