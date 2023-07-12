package equeue

import "encoding/json"

const KeyHeaderError = "x-errors-data"

const (
	MIMEJson     = "application/json"
	MIMEProtobuf = "application/x-protobuf"
)

type OutEventType string

const (
	OutEventTypeNormal  OutEventType = "normal"
	OutEventTypeFailure OutEventType = "failure"
	OutEventTypeDead    OutEventType = "dead"
)

// AsyncEvent :
type OutEvent struct {
	Status    OutEventType           // Repush when error (default is normal)
	ExpiredAt int64                  // Only use for failured status (delay(second) before push to queue)
	Queue     string                 // Queue name
	Exchange  string                 // routing key
	Data      interface{}            //
	Headers   map[string]interface{} //
}

type InEvent[T any] struct {
	Queue    string                 // Queue name
	Exchange string                 // routing key
	Data     *T                     //
	Headers  map[string]interface{} //
	Errors   *Error                 //
}

type Error struct {
	Count  int      `json:"count"`
	Errors []string `json:"errors"`
}

func (e Error) String() string {
	raw, _ := json.Marshal(&e)
	return string(raw)
}

func (e *Error) Add(err error) {
	e.Count++
	e.Errors = append(e.Errors, err.Error())
}

type FnConsumer[T any] func(*InEvent[T]) error

type ISerde[T any] interface {
	MIME() string

	Marshal(val *T) ([]byte, error)
	Unmarshal(raw []byte) (*T, error)
}

type ISerdeAny interface {
	MIME() string

	Marshal(val interface{}) ([]byte, error)
	Unmarshal(raw []byte, val interface{}) error
}

type IPublisher interface {
	SetSerde(sd ISerdeAny)

	Publish(ev *OutEvent) error
	PublishForce(ev *OutEvent)

	// CreateQueue(name string)
	// BindQueue()
}

type IGetter[T any] interface {
	Subcribe() <-chan *InEvent[T]
	UnSubscribe()
}
