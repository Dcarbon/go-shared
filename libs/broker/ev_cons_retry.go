package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Dcarbon/go-shared/libs/rabbit"
	"github.com/Dcarbon/go-shared/libs/utils"
	"github.com/streadway/amqp"
)

const keyHeaderError = "x-errors-data"

type CBRetry func(*amqp.Delivery) error
type CBRetryHook func(name string, d *amqp.Delivery)

type ErrorList struct {
	Errors []string `json:"errors"`
}

func (e ErrorList) String() string {
	raw, _ := json.Marshal(&e)
	return string(raw)
}

type ConsumerRetry struct {
	name        string
	maxRetry    int
	rbConn      rabbit.IConnection
	rbChannel   rabbit.IChannel
	cancel      context.CancelFunc
	cb          CBRetry
	onDead      CBRetryHook
	timeoutBase int
}

func NewEventConsumerRetry(
	name string,
	rbConn rabbit.IConnection,
	cbHandle CBRetry,
) (*ConsumerRetry, error) {
	ch, err := rbConn.Channel()
	if nil != err {
		return nil, err
	}
	defer ch.Close()

	var cons = &ConsumerRetry{
		name:        name,
		maxRetry:    5,
		rbConn:      rbConn,
		cancel:      nil,
		cb:          cbHandle,
		timeoutBase: 60 * 1000,
	}

	err = ch.ExchangeDeclare(name, "direct", true, false, false, false, nil)
	if nil != err {
		return nil, err
	}

	// Execute queue
	_, err = ch.QueueDeclare(cons.getExecuteQueue(), true, false, false, false, nil)
	if nil != err {
		return nil, err
	}

	err = ch.QueueBind(cons.getExecuteQueue(), "", name, false, nil)
	if nil != err {
		return nil, err
	}

	// Timeout queue
	_, err = ch.QueueDeclare(
		cons.getTimeoutQueue(),
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-dead-letter-exchange":    name,
			"x-dead-letter-routing-key": "",
		})
	if nil != err {
		return nil, err
	}

	// Dead queue
	_, err = ch.QueueDeclare(cons.getDeadQueue(), true, false, false, false, nil)
	if nil != err {
		return nil, err
	}

	ch2, err := rbConn.Channel()
	if nil != err {
		return nil, err
	}
	cons.rbChannel = ch2

	return cons, nil
}

// Start : run with goroutine
func (cons *ConsumerRetry) Start() {
	go cons.consumeEvent()
}

func (cons *ConsumerRetry) Stop() {
	if cons.cancel != nil {
		cons.cancel()
	}

}

// Base is milisecond
func (cons *ConsumerRetry) SetTimeoutBase(i int) *ConsumerRetry {
	if i < 100 {
		panic("Timeout base must be gte 100")
	}
	cons.timeoutBase = i
	return cons
}

func (cons *ConsumerRetry) SetOnDead(hook CBRetryHook,
) *ConsumerRetry {
	cons.onDead = hook
	return cons
}

func (cons *ConsumerRetry) consumeEvent() {
	if nil != cons.cancel {
		cons.cancel()
	}

	ch, err := cons.rbConn.Channel()
	utils.PanicError("Create amqp channel error", err)
	defer ch.Close()

	err = ch.Qos(10, 0, false)
	if nil != err {
		log.Println(
			"Setup qos of channel for consumer error: ",
			err,
			cons.getExecuteQueue(),
		)
	}

	msgChan, err := ch.Consume(
		cons.getExecuteQueue(),
		"",
		true,
		false,
		false,
		false,
		amqp.Table{},
	)
	utils.PanicError("Creat consumer error", err)

	ctx := context.TODO()
	ctx, cons.cancel = context.WithCancel(ctx)

	var done = false
	for {
		if done {
			cons.cancel = nil
			break
		}

		select {
		case <-ctx.Done():
			log.Println("Stop event Consumer")
			done = true
		case d, ok := <-msgChan:
			if !ok {
				log.Println(
					"Event consumer was close from outside. ",
					cons.getExecuteQueue(),
				)
				cons.cancel()
				done = true
			} else {
				err := cons.cb(&d)
				if nil != err {
					log.Println("Error: ", err)
					cons.repush(&d, err)
				}
			}
		}
	}
}

func (cons *ConsumerRetry) repush(d *amqp.Delivery, err error) {
	var headers = amqp.Table{}
	var d2 = amqp.Publishing{
		Body:        d.Body,
		Headers:     headers,
		ContentType: "application/json",
	}
	var errs = GetErrors(d)
	errs.Errors = append(errs.Errors, err.Error())

	d2.Headers[keyHeaderError] = errs.String()
	var queue = cons.getTimeoutQueue()
	var szErr = len(errs.Errors)
	if szErr >= cons.maxRetry {
		log.Printf("Push event %s to dead queue:%s\n", string(d2.Body), queue)
		if cons.onDead != nil {
			cons.onDead(cons.name, d)
		}
	} else {
		d2.Expiration = fmt.Sprintf("%d", (szErr+1)*(szErr+1)/3*cons.timeoutBase)
	}

	cons.rbChannel.PublishForce("", queue, false, false, d2)
}

func (cons *ConsumerRetry) getExecuteQueue() string {
	return cons.name + "-execute"
}

func (cons *ConsumerRetry) getTimeoutQueue() string {
	return fmt.Sprintf("%s-timeout", cons.name)
}

func (cons *ConsumerRetry) getDeadQueue() string {
	return cons.name + "-dead"
}
