package ievent

import (
	"encoding/json"
	"log"
)

// AsyncEvent :
type Event struct {
	Queue    string      // Queue name
	Exchange string      // routing key
	Data     interface{} //
	Headers  map[string]interface{}
}

func (ev *Event) Encode() ([]byte, error) {
	var raw []byte
	var err error

	if tmp, ok := ev.Data.([]byte); ok {
		raw = tmp
	} else {
		raw, err = json.Marshal(ev.Data)
		if nil != err {
			log.Println("marshal event error: ", ev, err)
			return nil, err
		}
	}

	return raw, nil
}

type IPublisher interface {
	Publish(ev *Event) error
}
