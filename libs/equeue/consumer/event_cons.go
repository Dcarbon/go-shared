package consumer

import (
	"log"

	"github.com/Dcarbon/go-shared/libs/equeue"
)

type Consumer[T any] struct {
	publisher equeue.IPublisher
	getter    equeue.IGetter[T]
	callback  equeue.FnConsumer[T]
}

func NewEventConsumer[T any](publisher equeue.IPublisher, getter equeue.IGetter[T],
) *Consumer[T] {
	var cons = &Consumer[T]{
		publisher: publisher,
		getter:    getter,
	}
	return cons
}

// Start : run with goroutine
func (cons *Consumer[T]) Start() {
	go cons.consume()
}

func (cons *Consumer[T]) Stop() {
	cons.getter.UnSubscribe()
}

func (cons *Consumer[T]) consume() {
	var eventQueue = cons.getter.Subcribe()
	for ev := range eventQueue {
		if nil == ev {
			break
		}

		err := cons.callback(ev)
		if nil != err {
			log.Println("Handle event error: ", ev.Data, err)
		}
	}
}
