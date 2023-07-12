package broker

// AsyncEvent :
type Event struct {
	Queue    string      // Queue name
	Exchange string      // routing key
	Data     interface{} //
	Headers  map[string]interface{}
}

type IPublisher interface {
	Publish(ev Event) error
}
