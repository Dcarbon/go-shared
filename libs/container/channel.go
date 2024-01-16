package container

import (
	"context"
	"log"
)

type Channel[T any] struct {
	// src     Source[T]
	ch       <-chan T
	onEvent  func(T)
	onClosed func()
	cancel   context.CancelFunc
}

func NewChannel[T any](ch chan T, onEvent func(T),
) *Channel[T] {
	var cc = &Channel[T]{
		ch:      ch,
		onEvent: onEvent,
	}
	return cc
}

func (cc *Channel[T]) SetSource(ch chan T) {
	if cc.cancel != nil {
		log.Println("[Warning] Change when channel is running")
	}
	cc.ch = ch
}

func (cc *Channel[T]) Start(wait bool) {
	if wait {
		cc.run()
	} else {
		go cc.run()
	}
}

func (cc *Channel[T]) StartWait() {

}

func (cc *Channel[T]) Stop() {
	if nil != cc.cancel {
		cc.cancel()
		cc.cancel = nil
	}
}

func (cc *Channel[T]) run() {
	if cc.cancel != nil {
		cc.cancel()
	}
	var done = false
	var ctx, cancel = context.WithCancel(context.Background())
	cc.cancel = cancel

	for {
		if done {
			if cc.cancel != nil {
				cc.cancel()
				cc.cancel = nil
			}

			if cc.onClosed != nil {
				cc.onClosed()
			}
			break
		}

		select {
		case ev, ok := <-cc.ch:
			if !ok {
				done = true
			} else {
				cc.onEvent(ev)
			}
		case <-ctx.Done():
			done = true
		}
	}
}
