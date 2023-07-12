package broker

import (
	"context"
	"log"

	"github.com/Dcarbon/go-shared/libs/rabbit"
	"github.com/Dcarbon/go-shared/libs/utils"
	"github.com/streadway/amqp"
)

type ConsumerHandler func(*amqp.Delivery)

type EventConsumer struct {
	queueName string
	rbConn    rabbit.IConnection
	cb        ConsumerHandler
	cancel    context.CancelFunc
}

func NewEventConsumer(queueName string, cb ConsumerHandler, conn rabbit.IConnection,
) *EventConsumer {
	var consumer = &EventConsumer{
		queueName: queueName,
		rbConn:    conn,
		cb:        cb,
	}
	return consumer
}

func NewEventConsumerWithURL(queueName string, cb ConsumerHandler, connURL string,
) (*EventConsumer, error) {
	conn, err := rabbit.Dial(connURL)
	if nil != err {
		return nil, err
	}
	var consumer = &EventConsumer{
		queueName: queueName,
		rbConn:    conn,
		cb:        cb,
	}
	return consumer, nil
}

// Start : run with goroutine
func (consumer *EventConsumer) Start() {
	go consumer.consumeEvent()
}

func (consumer *EventConsumer) Stop() {
	if consumer.cancel != nil {
		consumer.cancel()
	}
}

func (consumer *EventConsumer) consumeEvent() {
	if nil != consumer.cancel {
		consumer.cancel()
	}

	ch, err := consumer.rbConn.Channel()
	utils.PanicError("Create amqp channel error", err)
	defer ch.Close()

	err = ch.Qos(10, 0, false)
	if nil != err {
		log.Println("Setup qos of channel for consumer error: ", err, consumer.queueName)
	}

	msgChan, err := ch.Consume(
		consumer.queueName,
		"",
		true,
		false,
		false,
		false,
		amqp.Table{},
	)
	utils.PanicError("Creat consumer error", err)

	ctx := context.TODO()
	ctx, consumer.cancel = context.WithCancel(ctx)

	var done = false
	for {
		if done {
			break
		}

		select {
		case <-ctx.Done():
			log.Println("Stop event Consumer")
			consumer.cancel()
			done = true
		case d, ok := <-msgChan:
			if !ok {
				log.Println("Event consumer was close from outside. ", consumer.queueName)
				consumer.cancel()
				done = true
			} else {
				consumer.cb(&d)
			}
		}
	}
}
