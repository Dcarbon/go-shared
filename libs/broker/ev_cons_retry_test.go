package broker

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/Dcarbon/go-shared/libs/rabbit"
	"github.com/streadway/amqp"
)

type ERetryTest struct {
	Rand int `json:"rand"`
}

func TestEventConsumerRetry(t *testing.T) {
	var amqpUrl = "amqp://rbuser:244466666@localhost"
	var rbConn, err = rabbit.Dial(amqpUrl)
	if nil != err {
		log.Fatalln("Dial to rabbit error: ", err)
	}
	var name = "test-cons-retry"

	cons, err := NewEventConsumerRetry(
		name,
		rbConn,
		consEventForConsumerRetry,
	)
	if nil != err {
		log.Fatalln("Create consumer retry: ", err)
	}

	cons.SetOnDead(hookOnDead).SetTimeoutBase(1 * 1000)
	cons.Start()

	go pushEventForConsumerRetry(rbConn, name)
	time.Sleep(50 * time.Minute)
}

func pushEventForConsumerRetry(rbConn rabbit.IConnection, x string) {
	var ch, err = rbConn.Channel()
	if nil != err {
		log.Fatalln("")
	}

	rand.Seed(time.Now().Unix())
	for i := 0; i < 30; i++ {
		ch.Publish(x, "", false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(fmt.Sprintf(`{"rand": %d}`, rand.Int31n(10))),
		})

		time.Sleep(5 * time.Second)

	}

	// ch.Publish(x, "", false, false, amqp.Publishing{
	// 	ContentType: "application/json",
	// 	Body:        []byte(`{"rand": 1}`),
	// })
	// ch.Publish(x, "", false, false, amqp.Publishing{
	// 	ContentType: "application/json",
	// 	Body:        []byte(`{"rand": 2}`),
	// })
	// ch.Publish(x, "", false, false, amqp.Publishing{
	// 	ContentType: "application/json",
	// 	Body:        []byte(`{"rand": 3}`),
	// })

}

func consEventForConsumerRetry(d *amqp.Delivery) error {
	var ev = &ERetryTest{}
	var err = json.Unmarshal(d.Body, ev)
	if nil != err {
		return err
	}
	var el = GetErrors(d)
	log.Println()
	log.Println("Handle event with value ", ev.Rand, el.String())

	var m = ev.Rand % 3
	if m == 0 {
		log.Println("Handle event with value success", ev.Rand, el)
		return nil
	}
	if m == 1 {
		if len(el.Errors) >= 2 {
			log.Println("Handle event with value success", ev.Rand, el)
			return nil
		}
		return fmt.Errorf("failed by %d", ev.Rand)
	}

	return fmt.Errorf("always failed by %d", len(el.Errors))
}

func hookOnDead(name string, d *amqp.Delivery) {
	log.Println("Hook on dead: ", name, string(d.Body))
}
