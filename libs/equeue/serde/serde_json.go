package serde

import (
	"encoding/json"

	"github.com/Dcarbon/go-shared/libs/equeue"
)

type jsonSerde[T any] struct {
}

func NewJsonSerde[T any]() equeue.ISerde[T] {
	var jsd = &jsonSerde[T]{}
	return jsd
}

func (jsd *jsonSerde[T]) MIME() string {
	return "application/json"
}

func (jsd *jsonSerde[T]) Marshal(val *T) ([]byte, error) {
	return json.Marshal(val)
}

func (jsd *jsonSerde[T]) Unmarshal(raw []byte) (*T, error) {
	var rs = new(T)
	var err = json.Unmarshal(raw, rs)
	if nil != err {
		return nil, err
	}

	return rs, nil
}

type jsonSerdeAny struct {
}

func NewJsonSerdeAny() equeue.ISerdeAny {
	var jsd = &jsonSerdeAny{}
	return jsd
}

func (jsd *jsonSerdeAny) MIME() string {
	return "application/json"
}

func (jsd *jsonSerdeAny) Marshal(val interface{}) ([]byte, error) {
	return json.Marshal(val)
}

func (jsd *jsonSerdeAny) Unmarshal(raw []byte, val interface{}) error {
	return json.Unmarshal(raw, val)
}
