package rcons

import (
	"log"

	"github.com/Dcarbon/go-shared/libs/equeue"
)

const (
	delayQueue = "_delay_queue"
	deadQueue  = "_dead_queue"
)

type ConsumerRetry[T any] struct {
	timeoutBase int
	maxRetry    int
	publisher   equeue.IPublisher
	getter      equeue.IGetter[T]
	callback    equeue.FnConsumer[T]
}

func NewEventConsumerRetry[T any](publisher equeue.IPublisher, getter equeue.IGetter[T],
) (*ConsumerRetry[T], error) {
	var cons = &ConsumerRetry[T]{
		maxRetry:    5,
		publisher:   publisher,
		getter:      getter,
		timeoutBase: 60 * 1000,
	}
	return cons, nil
}

func (cons *ConsumerRetry[T]) SetTimeoutBase(i int) *ConsumerRetry[T] {
	if i < 100 {
		panic("Timeout base must be gte 100")
	}
	cons.timeoutBase = i
	return cons
}

// Start : run with goroutine
func (cons *ConsumerRetry[T]) Start() {
	cons.Stop()
	go cons.consume()
}

func (cons *ConsumerRetry[T]) Stop() {
	if cons.getter != nil {
		cons.getter.UnSubscribe()
	}
}

func (cons *ConsumerRetry[T]) consume() {
	var eventQueue = cons.getter.Subcribe()

	for ev := range eventQueue {
		if nil == ev {
			break
		}

		err := cons.callback(ev)
		if nil != err {
			log.Println("Handle event error: ", ev.Data, err)
			cons.repush(ev, err)
		}
	}
}

func (cons *ConsumerRetry[T]) repush(ev *equeue.InEvent[T], err error) {
	ev.Errors.Add(err)
	var outEvent = &equeue.OutEvent{
		Status:   equeue.OutEventTypeDead,
		Queue:    ev.Queue,
		Exchange: ev.Exchange,
		Data:     ev.Data,
		Headers: map[string]interface{}{
			equeue.KeyHeaderError: ev.Errors.String(),
		},
	}
	cons.publisher.Publish(outEvent)
}

// func (cons *ConsumerRetry) getExecuteQueue() string {
// 	return cons.name + "-execute"
// }

// func (cons *ConsumerRetry) getTimeoutQueue() string {
// 	return fmt.Sprintf("%s-timeout", cons.name)
// }

// func (cons *ConsumerRetry) getDeadQueue() string {
// 	return cons.name + "-dead"
// }
