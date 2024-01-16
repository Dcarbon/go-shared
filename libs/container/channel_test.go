package container

import (
	"log"
	"testing"
	"time"
)

func TestChannel(t *testing.T) {
	var arr = []int{}
	var ch = make(chan int, 128)
	var onEvent = func(i int) {
		arr = append(arr, i)
	}
	var push = func(v int) {
		ch <- v
	}

	go push(1)

	arr = append(arr, 2)
	var cc = NewChannel[int](ch, onEvent)
	cc.Start(false)
	go push(3)
	go push(4)
	go push(5)
	time.Sleep(1 * time.Second)
	cc.Stop()

	go push(10)

	time.Sleep(3 * time.Second)
	log.Println("Array: ", arr)
}
