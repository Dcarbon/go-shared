package rabbit

import "github.com/streadway/amqp"

type IConnection interface {
	Channel() (IChannel, error)
	IsClose() bool
	Close() bool
}

type IChannel interface {
	IsClose() bool

	Close() error
	Qos(prefetchCount, prefetchSize int, global bool) error
	Cancel(consumer string, noWait bool) error

	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)

	PublishForce(exchange, key string, mandatory, immediate bool, msg amqp.Publishing)
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	Ack(tag uint64, multiple bool) error
	Nack(tag uint64, multiple bool, requeue bool) error

	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	QueueDeclarePassive(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	QueueInspect(name string) (amqp.Queue, error)
	QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error
	QueueUnbind(name, key, exchange string, args amqp.Table) error
	QueuePurge(name string, noWait bool) (int, error)
	QueueDelete(name string, ifUnused, ifEmpty, noWait bool) (int, error)

	ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
	ExchangeDeclarePassive(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
	ExchangeDelete(name string, ifUnused, noWait bool) error
	ExchangeBind(destination, key, source string, noWait bool, args amqp.Table) error
	ExchangeUnbind(destination, key, source string, noWait bool, args amqp.Table) error
}
