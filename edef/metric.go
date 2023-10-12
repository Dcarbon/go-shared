package edef

import (
	"github.com/Dcarbon/go-shared/dmodels"
	"github.com/Dcarbon/go-shared/libs/ievent"
)

const (
	XSensorMetric = "sensor_metric"
)

// Sensor metric signature
type SMSign struct {
	IsIotSign  bool               `json:"isIotSign" ` //
	IotID      int64              `json:"iotId"`      //
	SensorID   int64              `json:"sensorId" `  //
	SensorType dmodels.SensorType `json:"sensorType"` //
	Data       string             `json:"data" `      // Hex json of SensorMetricExtract
	Signed     string             `json:"signed"`     // RSV Data
	Signer     dmodels.EthAddress `json:"signer"`     //
}

// SensorPusher: Sensor event pusher
type SensorPusher struct {
	pusher ievent.IPublisher
}

func NewSensorPusher(pusher ievent.IPublisher) *SensorPusher {
	var spusher = &SensorPusher{
		pusher: pusher,
	}
	return spusher
}

func (spusher *SensorPusher) PushNewMetric(sign *SMSign) {
	spusher.pusher.Publish(&ievent.Event{
		Exchange: XSensorMetric,
		Data:     sign,
	})
}
