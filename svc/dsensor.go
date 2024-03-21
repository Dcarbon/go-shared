package svc

import (
	"context"
	"time"

	"github.com/Dcarbon/arch-proto/pb"
	"github.com/Dcarbon/go-shared/gutils"
)

type ISensorService interface {
	GetById(ctx context.Context, req *pb.IdInt64) (*pb.Sensor, error)
}

type sensorService struct {
	client pb.SensorServiceClient
}

func (s sensorService) GetById(ctx context.Context, req *pb.IdInt64) (*pb.Sensor, error) {
	sensor, err := s.client.GetById(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.Sensor{
		Id:      sensor.Id,
		IotId:   sensor.IotId,
		Address: sensor.Address,
		Status:  sensor.Status,
		Type:    sensor.Type,
	}, nil
}

func NewSensorService(host string) (ISensorService, error) {
	cc, err := gutils.GetCCTimeout(host, 5*time.Second)
	if nil != err {
		return nil, err
	}

	var client = &sensorService{
		client: pb.NewSensorServiceClient(cc),
	}
	return client, nil
}
