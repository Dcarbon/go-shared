package rcons

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Dcarbon/go-shared/libs/equeue"
	"github.com/Dcarbon/go-shared/libs/rabbit"
	"github.com/streadway/amqp"
)

type Subscriber[T any] struct {
	queueName string
	autoAck   bool
	conn      rabbit.IConnection
	cancel    context.CancelFunc
	sd        equeue.ISerde[T]
	cb        equeue.FnConsumer[T]
}

func NewSubscriber[T any](conn rabbit.IConnection, qname string) *Subscriber[T] {
	var sub = &Subscriber[T]{
		queueName: qname,
		autoAck:   false,
		conn:      conn,
	}
	return sub
}

func (sub *Subscriber[T]) SetAutoAck(autoAck bool) {
	sub.autoAck = autoAck
}

func (sub *Subscriber[T]) SetCallback(cb equeue.FnConsumer[T]) {
	sub.cb = cb
}

func (sub *Subscriber[T]) Start() {
	go sub.consume()
}

func (sub *Subscriber[T]) Stop() {
	if sub.cancel != nil {
		sub.cancel()
	}
}

func (sub *Subscriber[T]) consume() error {
	ch, err := sub.conn.Channel()
	if nil != err {
		return fmt.Errorf("create amqp channel error: %s", err.Error())
	}
	defer ch.Close()

	err = ch.Qos(10, 0, false)
	if nil != err {
		log.Println("Setup qos of channel for consumer error: ", err, sub.queueName)
	}

	msgChan, err := ch.Consume(
		sub.queueName,
		"",
		true,
		false,
		false,
		false,
		amqp.Table{},
	)
	panicError("Creat consumer error", err)

	ctx := context.TODO()
	ctx, sub.cancel = context.WithCancel(ctx)

	var done = false
	for {
		if done {
			break
		}

		select {
		case <-ctx.Done():
			log.Println("Stop event Consumer")
			sub.cancel()
			done = true
		case d, ok := <-msgChan:
			if !ok {
				log.Println("Event consumer was close from outside. ", sub.queueName)
				sub.cancel()
				done = true
			} else {
				sub.handleEvent(&d)
			}
		}
	}

	return nil
}

func (sub *Subscriber[T]) handleEvent(d *amqp.Delivery) {
	var payload, err = sub.sd.Unmarshal(d.Body)
	if nil != err {
		log.Println("Marshal event error: ", err, string(d.Body), sub.queueName)
		d.Ack(false)
		return
	}

	var errHeader = GetErrors(d)
	var ev = &equeue.InEvent[T]{
		Queue:    sub.queueName,
		Exchange: "",
		Data:     payload,
		Headers:  d.Headers,
		Errors:   errHeader,
	}

	if !sub.autoAck {
		defer func() {
			e2 := recover()
			if nil != err {
				log.Println("Consumer crash", sub.queueName, e2)
				d.Nack(false, false)
				return
			}

			if err != nil {
				d.Nack(false, false)
			} else {
				d.Ack(false)
			}
		}()
	}
	err = sub.cb(ev)
}

func GetErrors(d *amqp.Delivery) *equeue.Error {
	var rawErr, ok = d.Headers[equeue.KeyHeaderError].(string)
	if !ok {
		rawErr = "{}"
	}

	var e = &equeue.Error{}
	json.Unmarshal([]byte(rawErr), e)
	return e
}
