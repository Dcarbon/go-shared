package rabbit

import (
	"context"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type autoChannel struct {
	*amqp.Channel
	conn       *connection
	consumers  []*consumerCache
	delayEvent []*publishMessage
	isClose    bool
}

func newChannel(conn *connection, rbChan *amqp.Channel) *autoChannel {
	var ch = &autoChannel{
		conn:    conn,
		Channel: rbChan,
		isClose: false,
	}
	ch.autoConnect()
	return ch
}

func (ch *autoChannel) IsClose() bool {
	return ch.isClose
}

func (ch *autoChannel) Cancel(consumer string, noWait bool) error {
	var err = ch.Channel.Cancel(consumer, noWait)
	if nil != err {
		return err
	}
	for i, cons := range ch.consumers {
		if cons.consumer == consumer {
			var remove = ch.consumers[i]
			ch.consumers = append(ch.consumers[:i], ch.consumers[i+1:]...)
			remove.flush()
			break
		}
	}
	return nil
}

func (ch *autoChannel) Close() error {
	if nil == ch.Channel {
		return nil
	}
	return ch.Channel.Close()
}

func (ch *autoChannel) Publish(
	exchange, key string,
	mandatory, immediate bool,
	msg amqp.Publishing,
) error {
	if ch.isClose {
		log.Println("Push to delay event")
		ch.delayEvent = append(ch.delayEvent, &publishMessage{
			exchange:  exchange,
			key:       key,
			mandatory: mandatory,
			immediate: immediate,
			msg:       msg,
		})
		return nil
	}
	return ch.Channel.Publish(exchange, key, mandatory, immediate, msg)
}

func (ch *autoChannel) PublishForce(
	exchange, key string,
	mandatory, immediate bool,
	msg amqp.Publishing,
) {
	if ch.isClose {
		log.Println("Push to delay event")
		ch.delayEvent = append(
			ch.delayEvent,
			&publishMessage{
				exchange:  exchange,
				key:       key,
				mandatory: mandatory,
				immediate: immediate,
				msg:       msg,
			})
	}
	go func() {
		for i := 0; i < 5; i++ {
			err := ch.Channel.Publish(exchange, key, mandatory, immediate, msg)
			if nil == err {
				break
			}
			log.Printf(
				"Publish event %s to x:%s key:%s error: %s\n",
				string(msg.Body), exchange, key, err.Error(),
			)
			time.Sleep(5 * time.Second)
		}
	}()
}

func (ch *autoChannel) Consume(
	queue, consumer string,
	autoAck, exclusive, noLocal, noWait bool, args amqp.Table,
) (<-chan amqp.Delivery, error) {
	var cons = &consumerCache{
		channel:   ch,
		queue:     queue,
		consumer:  consumer,
		autoAck:   autoAck,
		exclusive: exclusive,
		noLocal:   noLocal,
		noWait:    noWait,
		args:      args,
	}
	ch.consumers = append(ch.consumers, cons)
	err := ch.connectConsumer(cons)
	if nil != err {
		return nil, err
	}

	return cons.out, nil
}

func (ch *autoChannel) autoConnect() {
	if ch.Channel == nil {
		return
	}

	go func() {
		var c = make(chan *amqp.Error) // was close by amqp.Channel
		ch.Channel.NotifyClose(c)
		err := <-c
		ch.isClose = true

		if nil == err {
			// log.Println("Rabbitmq close channel normally")
			ch.flush()
			return
		}

		log.Println("Rabbitmq channel receive signal close: ", err)
		ch.conn.waitForOpen()
		for {
			rbChan, err := ch.conn.Connection.Channel()
			if nil != err {
				log.Println("Rabbitmq channel try reconnect error: ", err)
			} else {
				log.Println("Rabbitmq channel try reconnect success")
				ch.Channel = rbChan
				ch.isClose = false
				ch.autoConnect()
				ch.reconnectConsumer()
				ch.pushlishDelayEvent()
				break
			}
		}
	}()
}

func (ch *autoChannel) connectConsumer(cons *consumerCache) error {
	in, err := ch.Channel.Consume(
		cons.queue,
		cons.consumer,
		cons.autoAck,
		cons.exclusive,
		cons.noLocal,
		cons.noWait, cons.args,
	)
	if nil != err {
		return err
	}
	cons.startBridge()
	if cons.out == nil {
		cons.out = make(chan amqp.Delivery)
	}
	cons.in = in

	return nil
}

func (ch *autoChannel) reconnectConsumer() {
	for _, cons := range ch.consumers {
		err := ch.connectConsumer(cons)
		if nil != err {
			log.Println("Reconnect consumer error: ", err, cons.queue)
		} else {
			log.Println("Reconnect consumer success ", cons.queue)
		}
	}
}

func (ch *autoChannel) pushlishDelayEvent() {
	if ch.isClose {
		return
	}
	var delays = ch.delayEvent
	ch.delayEvent = nil
	for _, ev := range delays {
		err := ch.Channel.Publish(ev.exchange, ev.key, ev.mandatory, ev.immediate, ev.msg)
		if nil != err {
			log.Println("Publish delay event error: ", err)
		}
	}
}

func (ch *autoChannel) flush() {
	for _, cons := range ch.consumers {
		cons.flush()
	}
}

type consumerCache struct {
	ctx       context.Context
	cancel    context.CancelFunc
	channel   *autoChannel
	queue     string
	consumer  string
	autoAck   bool
	exclusive bool
	noLocal   bool
	noWait    bool
	args      amqp.Table
	in        <-chan amqp.Delivery
	out       chan amqp.Delivery
}

func (cons *consumerCache) startBridge() {
	cons.stopBridge()
	cons.ctx, cons.cancel = context.WithCancel(context.TODO())
	go func() {
		if nil == cons.in || nil == cons.out {
			log.Println("Start bridge with not enough in out queue: ", cons.queue)
			return
		}
		var done = false
		for {
			select {
			case msg, ok := <-cons.in:
				if !ok {
					done = true
				} else {
					cons.out <- msg
				}
			case <-cons.ctx.Done():
				done = true
			}

			if done {
				cons.stopBridge()
				break
			}
		}
	}()
}

func (cons *consumerCache) stopBridge() {
	if cons.cancel != nil {
		cons.cancel()
		cons.cancel = nil
	}
}

func (cons *consumerCache) flush() {
	if cons.out != nil {
		close(cons.out)
		cons.out = nil
	}
}

type publishMessage struct {
	exchange  string
	key       string
	mandatory bool
	immediate bool
	msg       amqp.Publishing
}
