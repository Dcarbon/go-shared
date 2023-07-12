package rabbit

import (
	"fmt"
	"log"
	"runtime"
	"testing"
	"time"

	"github.com/streadway/amqp"
)

const urlStr = "amqp://rbuser:2444@localhost/"
const queueName = "only-test"

func TestAll(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU() / 2)

	conn, err := Dial(urlStr)
	panicError("Connect to amqp ", err)

	ch, err := conn.Channel()
	panicError("Create channel ", err)

	ch.QueueDeclare(queueName, true, false, false, false, nil)
	go func() {
		ev, err := ch.Consume(queueName, "consumer-1", true, false, false, false, nil)
		panicError("Consume ", err)
		for msg := range ev {
			log.Println("Receive event (c1): ", string(msg.Body))
		}
	}()

	createQueue(conn, "abc-1")
	createQueue(conn, "abc-2")
	createQueue(conn, "abc-3")
	createQueue(conn, "abc-4")
	createQueue(conn, "abc-5")
	createQueue(conn, "abc-6")
	createQueue(conn, "abc-7")

	log.Println("PASSS")

	go func() {
		time.Sleep(10 * time.Second)

		ch2, err := conn.Channel()
		panicError("Create channel ", err)
		defer ch2.Close()

		for i := 0; i < 50; i++ {
			time.Sleep(10 * time.Second)
			err = ch2.Publish("", queueName, false, false, amqp.Publishing{
				Body: []byte(fmt.Sprintf("This is test from golang test %d", i)),
			})
			if nil != err {
				log.Println("Publish messsage error: ", err)
			}
		}

	}()

	// go func() {
	// 	time.Sleep(30 * time.Second)
	// 	var command = exec.Command("docker-compose", "restart", "rabbit")
	// 	err := command.Run()
	// 	if nil != err {
	// 		log.Println("Restart rabbit error: ", err)
	// 	}
	// }()

	var cClose = make(chan int)
	go func() {
		time.Sleep(150 * time.Minute)
		ch.Close()
		cClose <- 1
	}()

	<-cClose

	conn.Close()
	log.Println("Test done")
}

func createQueue(conn IConnection, queueName string) {
	ch2, err := conn.Channel()
	panicError("Create channel ", err)
	defer ch2.Close()
	ch2.QueueDeclare(queueName, true, false, false, false, nil)
}

func panicError(msg string, err error) {
	if nil != err {
		panic(msg + " error: " + err.Error())
	}
}
