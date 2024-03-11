package edef

import (
	"fmt"
	"github.com/Dcarbon/go-shared/dmodels"
	"github.com/Dcarbon/go-shared/libs/ievent"
	"log"
)

const (
	XSensorMetric        = "sensor_metric"
	XSensorMetricCreated = "sensor_created"
)

// Sensor metric signature
type SMSign struct {
	IsIotSign  bool               `json:"isIotSign" ` //
	IotId      int64              `json:"iotId"`      //
	SensorId   int64              `json:"sensorId" `  //
	SensorType dmodels.SensorType `json:"sensorType"` //
	Data       string             `json:"data" `      // Hex json of SensorMetricExtract
	Signed     string             `json:"signed"`     // RSV Data
	Signer     dmodels.EthAddress `json:"signer"`     //
}

// Sensor metric signature
type EventSenSorMetricCreated struct {
	IoTID    int64                `json:"id"`
	SensorID int64                `json:"sensorId"`
	Status   dmodels.DeviceStatus `json:"status"`
	Address  string               `json:"address"`
	Location *dmodels.Coord       `json:"location"` // Hex json of SensorMetricExtract
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
func (spusher *SensorPusher) PushNewMetricToMapIoTListener(sign *EventSenSorMetricCreated) {
	err := spusher.pusher.Publish(&ievent.Event{
		Exchange: XSensorMetric,
		Data:     sign,
	})
	if err != nil {
		log.Println(fmt.Sprintf("[SensorISM]: Emit event insert, publish event error: %s", err))
		return
	}
	log.Println(fmt.Sprintf("[SensorISM]: Emit event insert success. "))
}
