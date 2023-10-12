package ievent

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

type ErrorList struct {
	Errors []string `json:"errors"`
}

func (e ErrorList) String() string {
	raw, _ := json.Marshal(&e)
	return string(raw)
}

type CBHanle func(*amqp.Delivery) error
type CBHook func(name string, d *amqp.Delivery)

// type CROption struct {
// }

type ConsumerRetry struct {
	name         string
	maxRetry     int
	timeoutBase  int
	hasDeadQueue bool
	rbChannel    rabbit.IChannel
	cancel       context.CancelFunc
	cb           CBHanle
	onDead       CBHook
	onAckError   CBHook
	pusher       IPublisher
}

func NewConsumerRetry(name string, rbConn rabbit.IConnection, cb CBHanle,
) (*ConsumerRetry, error) {
	var cons = &ConsumerRetry{
		hasDeadQueue: true,
		name:         name,
		maxRetry:     5,
		timeoutBase:  60 * 1000, // 60 second
		cancel:       nil,
		cb:           cb,
	}

	err := cons.init(rbConn)
	if nil != err {
		return nil, err
	}

	cons.rbChannel, err = rbConn.Channel()
	if nil != err {
		return nil, err
	}

	cons.pusher = NewAsyncPusher(cons.rbChannel, 256)
	return cons, nil
}

func (cons *ConsumerRetry) GetMaxRetry() int {
	return cons.maxRetry
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

func (cons *ConsumerRetry) IsRunning() bool {
	return cons.cancel != nil
}

// Base is milisecond
func (cons *ConsumerRetry) SetTimeoutBase(i int) *ConsumerRetry {
	if i < 100 {
		panic("Timeout base must be gte 100")
	}
	cons.timeoutBase = i
	return cons
}

func (cons *ConsumerRetry) SetOnDead(hook CBHook,
) *ConsumerRetry {
	cons.onDead = hook
	return cons
}

func (cons *ConsumerRetry) SetOnAckError(hook CBHook,
) *ConsumerRetry {
	cons.onAckError = hook
	return cons
}

func (cons *ConsumerRetry) Purge() {
	cons.rbChannel.QueuePurge(cons.getExecuteQueue(), true)
	cons.rbChannel.QueuePurge(cons.GetTimeoutQueue(), true)
}

func (cons *ConsumerRetry) consumeEvent() {
	if nil != cons.cancel {
		cons.cancel()
	}

	err := cons.rbChannel.Qos(5, 0, false)
	if nil != err {
		log.Println(
			"Setup qos of channel for consumer error: ",
			err,
			cons.getExecuteQueue(),
		)
	}

	msgChan, err := cons.rbChannel.Consume(
		cons.getExecuteQueue(),
		"",
		false,
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
			// log.Println("Stop event consumer " + cons.name)
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
				err2 := d.Ack(false)
				if nil != err2 {
					log.Println("Ack message error: ", err2, string(d.Body))
				}
				if nil != err {
					cons.repush(&d, err)
				}
			}
		}
	}
}

func (cons *ConsumerRetry) repush(d *amqp.Delivery, eInput error) {
	var headers = amqp.Table{}
	var pubVal = amqp.Publishing{
		Body:        d.Body,
		Headers:     headers,
		ContentType: "application/json",
	}

	var errs = GetErrors(d)
	errs.Errors = append(errs.Errors, eInput.Error())
	if nil == d.Headers {
		d.Headers = amqp.Table{}
	}
	d.Headers[keyHeaderError] = errs.String()

	pubVal.Headers[keyHeaderError] = errs.String()
	var queue = ""
	var szErr = len(errs.Errors)
	if szErr >= cons.maxRetry {
		if cons.hasDeadQueue {
			queue = cons.getDeadQueue()
			log.Printf("Push event %s to dead queue:%s\n", string(pubVal.Body), queue)
		}

		if cons.onDead != nil {
			cons.onDead(cons.name, d)
		}

	} else {
		queue = cons.GetTimeoutQueue()
		pubVal.Expiration = fmt.Sprintf("%d", (szErr+1)*(szErr+1)/3*cons.timeoutBase)
	}

	if queue != "" {
		err := cons.rbChannel.Publish("", queue, false, false, pubVal)
		if nil != err {
			d.Nack(false, true)
		}
	}
}

func (cons *ConsumerRetry) init(rbConn rabbit.IConnection) error {
	ch, err := rbConn.Channel()
	if nil != err {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(cons.name, "direct", true, false, false, false, nil)
	if nil != err {
		return err
	}

	// Execute queue
	_, err = ch.QueueDeclare(cons.getExecuteQueue(), true, false, false, false, nil)
	if nil != err {
		return err
	}

	err = ch.QueueBind(cons.getExecuteQueue(), "", cons.name, false, nil)
	if nil != err {
		return err
	}

	// Timeout queue
	_, err = ch.QueueDeclare(
		cons.GetTimeoutQueue(),
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-dead-letter-exchange":    cons.name,
			"x-dead-letter-routing-key": "",
		})
	if nil != err {
		return err
	}

	if cons.hasDeadQueue {
		// Dead queue
		_, err = ch.QueueDeclare(
			cons.getDeadQueue(), true, false, false, false, nil,
		)
		if nil != err {
			return err
		}
	}

	return nil
}

func (cons *ConsumerRetry) getExecuteQueue() string {
	return cons.name + "-x"
}

func (cons *ConsumerRetry) GetTimeoutQueue() string {
	return fmt.Sprintf("%s-timeout", cons.name)
}

func (cons *ConsumerRetry) getDeadQueue() string {
	return cons.name + "-dead"
}

func (cons *ConsumerRetry) Publish(ev *Event) error {
	// if ev.Delay > 0 {
	// 	ev.Exchange = ""
	// 	ev.Queue = cons.getTimeoutQueue()
	// }
	return cons.pusher.Publish(ev)
}
