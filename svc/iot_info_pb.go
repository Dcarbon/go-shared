package svc

import (
	"context"
	"time"

	"github.com/Dcarbon/arch-proto/pb"
	"github.com/Dcarbon/go-shared/dmodels"
	"github.com/Dcarbon/go-shared/gutils"
)

// Service client

type pbIotClient struct {
	iiot pb.IotServiceClient
}

func NewIotService(host string) (IIotInfo, error) {
	cc, err := gutils.GetCCTimeout(host, 5*time.Second)
	if nil != err {
		return nil, err
	}

	var client = &pbIotClient{
		iiot: pb.NewIotServiceClient(cc),
	}
	return client, nil
}

func (client *pbIotClient) GetById(id int64) (*Iot, error) {
	data, err := client.iiot.GetIot(context.TODO(), &pb.IdInt64{
		Id: id,
	})
	if nil != err {
		return nil, err
	}
	return &Iot{
		Id:      data.Id,
		Project: data.Project,
		Address: data.Adddress,
		Status:  dmodels.DeviceStatus(data.Status),
	}, nil
}

func (client *pbIotClient) GetByAddress(addr string) (*Iot, error) {
	data, err := client.iiot.GetByAddress(context.TODO(), &pb.EAddress{
		Address: addr,
	})
	if nil != err {
		return nil, err
	}
	return &Iot{
		Id:      data.Id,
		Project: data.Project,
		Address: data.Adddress,
		Status:  dmodels.DeviceStatus(data.Status),
	}, nil
}