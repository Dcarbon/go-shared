package svc

import (
	"errors"
	"strings"

	"github.com/Dcarbon/go-shared/dmodels"
)

var DefaultMockIot = []*Iot{
	{Id: 292, Project: 1, Address: "0xe445517abb524002bb04c96f96abb87b8b19b53d", Status: dmodels.DeviceStatusSuccess},
	{Id: 2, Project: 2, Address: "0x19Adf96848504a06383b47aAA9BbBC6638E81afD", Status: dmodels.DeviceStatusSuccess},
	{Id: 3, Project: 3, Address: "0x56D2C9baB06e391b470365be694671c0F1dE30EC", Status: dmodels.DeviceStatusSuccess},
	{Id: 4, Project: 4, Address: "0x91f4F05882cE70A5d2eEbb459DaA96c866fBd5E7", Status: dmodels.DeviceStatusSuccess},
	{Id: 5, Project: 5, Address: "0x67ed13A2a07b7473A66908CFb9B36169C388917b", Status: dmodels.DeviceStatusSuccess},
	{Id: 6, Project: 6, Address: "0x2aa82ec74bB8b2964647C46452B2fb259Dc9b187", Status: dmodels.DeviceStatusSuccess},
	{Id: 7, Project: 7, Address: "0xf06eF2592ce2328782d1040790e828ba9c43F6Ef", Status: dmodels.DeviceStatusSuccess},
	{Id: 8, Project: 8, Address: "0x80Cc331b778a19FBA10c5A99B12BC94Ad02eE31E", Status: dmodels.DeviceStatusRegister},
	{Id: 9, Project: 9, Address: "0x971b57c8a974254F26197d5a362169a8aceD684a", Status: dmodels.DeviceStatusReject},
}

type mockIotClient struct {
	data     []*Iot
	mAddress map[string]*Iot
	mId      map[int64]*Iot
}

func NewMockIotClient(data ...*Iot) *mockIotClient {
	var mockIot = &mockIotClient{
		data:     data,
		mAddress: make(map[string]*Iot, len(data)),
		mId:      make(map[int64]*Iot, len(data)),
	}

	for _, it := range data {
		mockIot.mAddress[strings.ToLower(it.Address)] = it
		mockIot.mId[it.Id] = it
	}

	return mockIot
}

func (miot *mockIotClient) GetById(id int64) (*Iot, error) {
	var iot = miot.mId[id]
	if nil != iot {
		return &Iot{
			Id:      iot.Id,
			Project: iot.Project,
			Address: iot.Address,
			Status:  iot.Status,
		}, nil
	}
	return nil, errors.New("iot not found")
}

func (miot *mockIotClient) GetByAddress(addr string) (*Iot, error) {
	var iot = miot.mAddress[strings.ToLower(addr)]
	if nil != iot {
		return &Iot{
			Id:      iot.Id,
			Project: iot.Project,
			Address: iot.Address,
			Status:  iot.Status,
		}, nil
	}
	return nil, errors.New("iot not found")
}

func (miot *mockIotClient) GetIotsActivated() (*[]Iot, error) {
	// TODO: implement
	return nil, nil
}
