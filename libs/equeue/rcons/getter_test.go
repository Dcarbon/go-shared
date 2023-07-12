package rcons

import (
	"testing"

	"github.com/Dcarbon/go-shared/libs/rabbit"
)

var rbUrl = ""
var rbConn rabbit.IConnection
var rbChannel rabbit.IChannel

var x1 = "x_1"
var x2 = "x_2"

var xDelay = "x_delay"
var xDead = "x_dead"

func init() {
	var err error
	rbConn, err = rabbit.Dial(rbUrl)
	panicError("", err)

	rbChannel, err = rbConn.Channel()
	panicError("", err)

	err = rbChannel.ExchangeDeclare(x1, "", true, false, false, false, nil)
	panicError("", err)

	err = rbChannel.ExchangeDeclare(x2, "", true, false, false, false, nil)
	panicError("", err)

	err = rbChannel.ExchangeDeclare(xDead, "", true, false, false, false, nil)
	panicError("", err)

	err = rbChannel.ExchangeDeclare(xDelay, "", true, false, false, false, nil)
	panicError("", err)
}

// https://www.rabbitmq.com/dlx.html
func TestDLXQueue(t *testing.T) {

}
