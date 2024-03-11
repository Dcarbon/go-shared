package svc

import (
	"context"
	"github.com/Dcarbon/arch-proto/pb"
)

//var DefaultMockSensor = []*pb.Sensor{&pb.Sensor{
//	Id:      1,
//	IotId:   1,
//	Address: "0xE445517AbB524002Bb04C96F96aBb87b8B19b53d",
//	Status:  1,
//	Type:    1,
//}}

type MockSensorClient struct {
}

func NewMockSensorService() *MockSensorClient {

	return &MockSensorClient{}
}

func (miot *MockSensorClient) GetById(ctx context.Context, req *pb.IdInt64) (*pb.Sensor, error) {
	return &pb.Sensor{
		Id:      1,
		IotId:   1,
		Address: "0xE445517AbB524002Bb04C96F96aBb87b8B19b53d",
		Status:  1,
		Type:    1,
	}, nil
}
