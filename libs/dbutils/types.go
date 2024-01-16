package dbutils

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"log"

	"gorm.io/gorm"
)

type MapSS map[string]string

func (m *MapSS) Scan(value interface{}) error {
	if nil == m {
		m = new(MapSS)
	}
	switch vt := value.(type) {
	case string:
		return json.Unmarshal([]byte(vt), m)
	case []byte:
		return json.Unmarshal(vt, m)
	}
	return errors.New("scan value type for MapSS invalid")
}

func (m MapSS) Value() (driver.Value, error) {
	if nil == m {
		return nil, nil
	}
	return json.Marshal(m)
}

type Strings []string

func (ss *Strings) Scan(value interface{}) error {
	if nil == ss {
		ss = new(Strings)
	}
	switch vt := value.(type) {
	case string:
		return json.Unmarshal([]byte(vt), ss)
	case []byte:
		return json.Unmarshal(vt, ss)
	}
	return errors.New("scan value type for Strings is invalid")
}

func (ss Strings) Value() (driver.Value, error) {
	if nil == ss {
		return nil, nil
	}
	return json.Marshal(ss)
}

type Transform[TOld any, TNew any] struct {
	TblOld    string
	TblNew    string
	OnConvert func(*TOld) (*TNew, error)
	OnLoaded  func(*gorm.DB, *TNew) error
}

func (trans *Transform[T1, T2]) LoadAll(dbOld, dbNew *gorm.DB) error {
	var skip = 0
	var limit = 100
	for {

		data := make([]*T1, 0, limit)
		err := dbOld.Table(trans.TblOld).Offset(skip).Limit(limit).Order("id asc").Find(&data).Error
		if nil != err {
			return err
		}

		for _, it := range data {
			converted, err := trans.OnConvert(it)
			if nil != err {

				log.Println()
				continue
			}

			err = trans.OnLoaded(dbNew, converted)
			if nil != err {
				log.Println("Transform error: ", err)
				// return err
			}
		}

		if len(data) != limit {
			break
		}
		skip += limit
	}
	return nil
}

func (trans *Transform[TOld, TNew]) Insert(dbNew *gorm.DB, data *TNew) error {
	return dbNew.Table(trans.TblNew).Create(data).Error
}
