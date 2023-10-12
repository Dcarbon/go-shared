package ievent

import (
	"encoding/json"

	"github.com/Dcarbon/go-shared/libs/rabbit"
	"github.com/Dcarbon/go-shared/libs/utils"
	"github.com/streadway/amqp"
)

// CreateFanoutExchange :
func CreateFanoutExchange(conn rabbit.IConnection, name string) {
	ch, err := conn.Channel()
	utils.PanicError("Create rabbitmq channel error", err)
	defer ch.Close()

	err = ch.ExchangeDeclare(name, "fanout", true, false, false, false, nil)
	utils.PanicError("Create fanout exchange error", err)
}

// CreateFanoutExchange :
func CreateDirectExchange(conn rabbit.IConnection, name, binding string) {
	ch, err := conn.Channel()
	utils.PanicError("Create rabbitmq channel error", err)
	defer ch.Close()

	err = ch.ExchangeDeclare(name, "direct", true, false, false, false, nil)
	utils.PanicError("Create direct exchange error", err)

	if binding != "" {
		ch.ExchangeBind(name, "", binding, false, nil)
	}
}

// CreateQueue :
func CreateQueue(conn rabbit.IConnection, name string, exchange string, headers map[string]interface{}) {
	ch, err := conn.Channel()
	utils.PanicError("Create rabbitmq channel error", err)
	defer ch.Close()

	_, err = ch.QueueDeclare(name, true, false, false, false, headers)
	utils.PanicError("Create queue error", err)

	if exchange != "" {
		err = ch.QueueBind(name, "", exchange, false, nil)
		utils.PanicError("Create queue error", err)
	}
}

func GetErrors(d *amqp.Delivery) *ErrorList {
	var rawErr, ok = d.Headers[keyHeaderError].(string)
	if !ok {
		rawErr = "{}"
	}
	var errs = &ErrorList{}
	json.Unmarshal([]byte(rawErr), errs)
	return errs
}
