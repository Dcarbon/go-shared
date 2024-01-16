package dmodels

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/Dcarbon/go-shared/ecodes"
)

var regString = regexp.MustCompile(`"*"$`)

type SensorType int32

const (
	SensorTypeNone        SensorType = 0
	SensorTypeFlow        SensorType = 1
	SensorTypePower       SensorType = 2
	SensorTypeGPS         SensorType = 3
	SensorTypeThermometer SensorType = 4
)

type DeviceStatus int32

const (
	DeviceStatusReject   DeviceStatus = 1
	DeviceStatusRegister DeviceStatus = 5
	DeviceStatusSuccess  DeviceStatus = 10
)

type Sort int

const (
	SortASC  Sort = 0
	SortDesc Sort = 1
)

func (s *Sort) String() string {
	if *s == SortASC {
		return "asc"
	}
	return "desc"
}

type DInterval int

const (
	DINone  DInterval = 0
	DIHour  DInterval = 1
	DIDay   DInterval = 2
	DIMonth DInterval = 3
	DIYear  DInterval = 4
)

func (di *DInterval) String() string {
	switch *di {
	case DINone:
		return ""
	case DIHour:
		return "hour"
	case DIDay:
		return "day"
	case DIMonth:
		return "month"
	case DIYear:
		return "yearn"
	default:
		return ""
	}
}

type DefaultMetric struct {
	Val Float64 `json:"value,omitempty" cql:"value"`
	// Params []Float64 `json:"params,omitempty" cql:"value"`
}

type GPSMetric struct {
	Lat Float64 `json:"lat,omitempty" cql:"lat"`
	Lng Float64 `json:"lng,omitempty" cql:"lng"`
}

type AllMetric struct {
	DefaultMetric
	GPSMetric
}

func (am *AllMetric) IsValid(sType SensorType) error {
	switch sType {
	case SensorTypeNone:
	case SensorTypeFlow:
		if am.DefaultMetric.Val <= 0 {
			return NewError(ecodes.SensorInvalidMetric, "Indicator of metric (flow) must be > 0")
		}
	case SensorTypePower:
		if am.DefaultMetric.Val <= 0 {
			return NewError(ecodes.SensorInvalidMetric, "Indicator of metric (power) must be > 0")
		}
	case SensorTypeThermometer:
		if am.Val < -10 && am.Val > 1000 {
			return NewError(ecodes.SensorInvalidMetric, "Indicator of metric (thermometer) must be in range [-10;1000]")
		}
	case SensorTypeGPS:
		if am.Lat == 0 && am.Lng == 0 {
			return NewError(ecodes.SensorInvalidMetric, "Indicator of metric (gps) must be != 0")
		}
	default:
		return NewError(ecodes.SensorInvalidType, "Sensor type is not existed")
	}
	return nil
}

func (am *AllMetric) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var rs = &AllMetric{}
	var err error
	switch vt := value.(type) {
	case string:
		if vt == `""` {
			return nil
		}
		err = json.Unmarshal([]byte(vt), rs)
	case []byte:
		err = json.Unmarshal(vt, rs)
	default:
		return errors.New("can't scan metric")
	}
	if nil != err {
		return err
	}
	if nil == am {
		am = new(AllMetric)
	}
	*am = *rs
	return nil
}

func (am AllMetric) Value() (driver.Value, error) {
	v, err := json.Marshal(am)
	if nil != err {
		log.Println("All value: ", am)
	}
	return v, err
}

type Float64 float64

func (f *Float64) MarshalJSON() ([]byte, error) {
	if nil == f {
		return []byte(""), nil
	}
	return []byte(fmt.Sprintf(`%f`, *f)), nil
}

func (f *Float64) UnmarshalJSON(data []byte) error {
	var s = string(data)
	if regString.Match(data) {
		s = s[1 : len(s)-1]
	}

	if strings.ToLower(s) == "nan" {
		// return errors.New("NaN is not allow")
		if nil == f {
			f = new(Float64)
		}
		*f = 0.0
	} else {
		v, err := strconv.ParseFloat(s, 64)
		if nil != err {
			return err
		}

		if nil == f {
			f = new(Float64)
		}

		*f = Float64(v)
	}

	return nil
}
