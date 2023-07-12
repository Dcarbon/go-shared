package rcons

import (
	"log"
	"time"

	"github.com/Dcarbon/go-shared/libs/equeue"
	"github.com/Dcarbon/go-shared/libs/equeue/serde"
	"github.com/Dcarbon/go-shared/libs/rabbit"
	"github.com/streadway/amqp"
)

// directPusher :
type DirectPusher struct {
	rbChan rabbit.IChannel
	sd     equeue.ISerdeAny
}

func NewDirectPusher(rbChan rabbit.IChannel) equeue.IPublisher {
	var ap = &DirectPusher{
		rbChan: rbChan,
		sd:     serde.NewJsonSerdeAny(),
	}
	return ap
}

func (ap *DirectPusher) SetSerde(sd equeue.ISerdeAny) {
	ap.sd = sd
}

func (ap *DirectPusher) Publish(ev *equeue.OutEvent) error {
	if ev.Status == "" {
		ev.Status = equeue.OutEventTypeNormal
	}
	switch ev.Status {
	case equeue.OutEventTypeNormal:
		return ap.push(ev)
	case equeue.OutEventTypeFailure:
		var now = time.Now().Unix()
		if ev.ExpiredAt < now {
			ev.ExpiredAt = now + 60
		}
		ap.PublishForce(ev)
		return nil
	case equeue.OutEventTypeDead:
		ap.PublishForce(ev)
		return nil
	default:
		return ap.push(ev)
	}
	// return fmt.Errorf("invalid event status: %+v", ev)
}

func (ap *DirectPusher) PublishForce(ev *equeue.OutEvent) {
	ap.push(ev)
}

func (ap *DirectPusher) push(ev *equeue.OutEvent) error {
	var raw []byte
	var err error

	if tmp, ok := ev.Data.([]byte); ok {
		raw = tmp
	} else {
		raw, err = ap.sd.Marshal(ev.Data)
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
			ContentType: ap.sd.MIME(),
			Body:        raw,
			Headers:     ev.Headers,
		})
	if nil != err {
		log.Printf("Push event to queue %s error: %s\n", ev.Queue, err.Error())
	}

	return err
}
