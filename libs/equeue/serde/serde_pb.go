package serde

import (
	"errors"

	"github.com/Dcarbon/go-shared/libs/equeue"
	"google.golang.org/protobuf/proto"
)

type pbSerde[T any] struct {
}

func NewPbSerde[T any]() equeue.ISerde[T] {
	var rs = new(T)
	var _, ok = interface{}(rs).(proto.Message)
	if !ok {
		panic("Invalid input type")
	}

	var jsd = &pbSerde[T]{}

	return jsd
}

func (psd *pbSerde[T]) MIME() string {
	return "application/x-protobuf"
}

func (psd *pbSerde[T]) Marshal(val *T) ([]byte, error) {
	var v, ok = interface{}(val).(proto.Message)
	if !ok {
		return nil, errors.New("value must be proto.Message")
	}

	return proto.Marshal(v)
}

func (psd *pbSerde[T]) Unmarshal(raw []byte) (*T, error) {
	var rs = new(T)
	var rsp, _ = interface{}(rs).(proto.Message)

	var err = proto.Unmarshal(raw, rsp)
	if nil != err {
		return nil, err
	}

	return rs, nil
}

// func (jsd *pbSerde[T]) Marshal2(val interface{}) ([]byte, error) {
// 	var v, ok = interface{}(val).(proto.Message)
// 	if !ok {
// 		return nil, errors.New("value must be proto.Message")
// 	}

// 	return proto.Marshal(v)
// }

// func (jsd *pbSerde[T]) Unmarshal2(raw []byte, val interface{}) error {
// 	var rsp, ok = val.(proto.Message)
// 	if !ok {
// 		return errors.New("type of value must be proto.Message")
// 	}

// 	var err = proto.Unmarshal(raw, rsp)
// 	if nil != err {
// 		return err
// 	}

// 	return nil
// }

type pbSerdeAny struct{}

func NewPbSerdeAny() equeue.ISerdeAny {
	var pa = &pbSerdeAny{}
	return pa
}

func (psd *pbSerdeAny) MIME() string {
	return "application/x-protobuf"
}

func (pa *pbSerdeAny) Marshal(val interface{}) ([]byte, error) {
	var v, ok = interface{}(val).(proto.Message)
	if !ok {
		return nil, errors.New("value must be proto.Message")
	}

	return proto.Marshal(v)
}

func (pa *pbSerdeAny) Unmarshal(raw []byte, val interface{}) error {
	var rsp, ok = val.(proto.Message)
	if !ok {
		return errors.New("type of value must be proto.Message")
	}

	var err = proto.Unmarshal(raw, rsp)
	if nil != err {
		return err
	}

	return nil
}
