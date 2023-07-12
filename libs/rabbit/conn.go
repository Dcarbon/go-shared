package rabbit

import (
	"errors"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type connection struct {
	url     string
	isClose bool
	*amqp.Connection
}

func Dial(urlStr string) (IConnection, error) {
	if urlStr == "" {
		return nil, errors.New("AMQP url is empty")
	}

	rbConn, err := amqp.Dial(urlStr)
	if nil != err {
		return nil, err
	}
	var rs = &connection{
		url:        urlStr,
		isClose:    false,
		Connection: rbConn,
	}
	rs.autoConnect()
	return rs, nil
}

func (conn *connection) autoConnect() {
	log.Println("Rabbitmq connection auto connect")
	if conn.Connection == nil {
		return
	}
	go func() {
		var c = make(chan *amqp.Error)
		conn.Connection.NotifyClose(c)
		err := <-c
		conn.isClose = true

		log.Println("Trying reconnect when error: ", err)

		conn.Close()
		for {
			rbConn, err := amqp.Dial(conn.url)
			if nil != err {
				log.Println("Connect to rabbitmq error: ", err)
				time.Sleep(5 * time.Second)
			} else {
				conn.isClose = false
				conn.Connection = rbConn
				conn.autoConnect()
				break
			}
		}

	}()
}

func (conn *connection) Channel() (IChannel, error) {
	rbChan, err := conn.Connection.Channel()
	if nil != err {
		return nil, err
	}
	return newChannel(conn, rbChan), nil
}

func (conn *connection) IsClose() bool {
	if nil == conn.Connection {
		return true
	}
	return conn.Connection.IsClosed()
}

func (conn *connection) Close() bool {
	return conn.isClose
}

func (conn *connection) waitForOpen() {
	for {
		if conn.isClose {
			time.Sleep(6 * time.Second)
		} else {
			break
		}
	}
}
