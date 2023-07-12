package rcons

import (
	"context"
	"fmt"
	"log"

	"github.com/Dcarbon/go-shared/libs/equeue"
	"github.com/Dcarbon/go-shared/libs/rabbit"
	"github.com/streadway/amqp"
)

type Getter[T any] struct {
	queueName string
	conn      rabbit.IConnection
	cancel    context.CancelFunc
	sd        equeue.ISerde[T]
	ch        chan *equeue.InEvent[T]
}

func NewGetter[T any](conn rabbit.IConnection, qname string) *Getter[T] {
	var g = &Getter[T]{
		queueName: qname,
		conn:      conn,
	}
	return g
}

func (g *Getter[T]) Subcribe() <-chan *equeue.InEvent[T] {
	if g.cancel != nil {
		panic("Only subscribe one time")
	}

	g.ch = make(chan *equeue.InEvent[T])
	go g.consume()
	return g.ch
}

func (g *Getter[T]) UnSubscribe() {
	if nil != g.cancel {
		close(g.ch)
		g.cancel()
		g.cancel = nil
		g.ch = nil
	}
}

func (g *Getter[T]) consume() error {
	ch, err := g.conn.Channel()
	if nil != err {
		return fmt.Errorf("create amqp channel error: %s", err.Error())
	}
	defer ch.Close()

	err = ch.Qos(10, 0, false)
	if nil != err {
		log.Println("Setup qos of channel for consumer error: ", err, g.queueName)
	}

	msgChan, err := ch.Consume(
		g.queueName,
		"",
		true,
		false,
		false,
		false,
		amqp.Table{},
	)
	panicError("Creat consumer error", err)

	ctx := context.TODO()
	ctx, g.cancel = context.WithCancel(ctx)

	var done = false
	for {
		if done {
			break
		}

		select {
		case <-ctx.Done():
			log.Println("Stop event Consumer")
			g.cancel()
			done = true
		case d, ok := <-msgChan:
			if !ok {
				log.Println("Event consumer was close from outside. ", g.queueName)
				g.cancel()
				done = true
			} else {
				g.handleEvent(&d)
			}
		}
	}

	return nil
}

func (g *Getter[T]) handleEvent(d *amqp.Delivery) {
	var payload, err = g.sd.Unmarshal(d.Body)
	if nil != err {
		log.Println("Marshal event error: ", err, string(d.Body), g.queueName)
		d.Ack(false)
		return
	}

	var errHeader = GetErrors(d)
	var ev = &equeue.InEvent[T]{
		Queue:    g.queueName,
		Exchange: "",
		Data:     payload,
		Headers:  d.Headers,
		Errors:   errHeader,
	}

	g.ch <- ev
}
