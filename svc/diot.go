package svc

import "github.com/Dcarbon/go-shared/dmodels"

type Iot struct {
	Id         int64                `json:"id"         `
	Project    int64                `json:"project"    `
	Address    string               `json:"address"    `
	Status     dmodels.DeviceStatus `json:"status"     `
	TimeRemain int64
}

type RemainTime struct {
	Id         int64 `json:"id"         `
	RemainTime int64 `json:"project"    `
}

type IIotInfo interface {
	GetById(id int64) (*Iot, error)
	GetByAddress(addr string) (*Iot, error)
	GetIotsActivated() (*[]Iot, error)
	UpdateRemainTime(req RemainTime) error
}
