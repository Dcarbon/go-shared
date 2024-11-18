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
		Id:         data.Id,
		Project:    data.Project,
		Address:    data.Address,
		Status:     dmodels.DeviceStatus(data.Status),
		TimeRemain: data.RemainTime,
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
		Address: data.Address,
		Status:  dmodels.DeviceStatus(data.Status),
	}, nil
}

func (client *pbIotClient) GetIotsActivated() (*[]Iot, error) {
	iots, err := client.iiot.GetIots(context.TODO(), &pb.RIotGetList{Status: int32(pb.IOTStatus_IOTS_Success.Number()), Type: -1})
	if nil != err {
		return nil, err
	}
	var result []Iot
	for _, data := range iots.Data {
		result = append(result, Iot{
			Id:      data.Id,
			Project: data.Project,
			Address: data.Address,
			Status:  dmodels.DeviceStatus(data.Status)})
	}
	return &result, nil
}

func (client *pbIotClient) UpdateRemainTime(req RemainTime) error {
	_, err := client.iiot.UpdateRemain(context.TODO(), &pb.RUpdateRemainTime{
		IotId:      req.Id,
		RemainTime: req.RemainTime,
	})
	if nil != err {
		return err
	}
	return nil
}
