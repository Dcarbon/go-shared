package dmodels

import (
	"database/sql/driver"
	"encoding/hex"
	"errors"
	"fmt"
	"math"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
)

// Coord:
type Coord struct {
	SRID int     `json:"srid"` //
	Lat  float64 `json:"lat"`  // vi tuyen (pgis: y)
	Lng  float64 `json:"lng"`  // kinh tuyen:(pgis: x)
}

func NewCoord(srid int, lng, lat float64) *Coord {
	return &Coord{
		SRID: srid,
		Lat:  lat,
		Lng:  lng,
	}
}

func NewCoord4326(lng, lat float64) *Coord {
	return &Coord{
		SRID: 4326,
		Lat:  lat,
		Lng:  lng,
	}
}

func NewCoord3857(lng, lat float64) *Coord {
	return &Coord{
		SRID: 3857,
		Lat:  lat,
		Lng:  lng,
	}
}

func (p *Coord) String() string {
	return fmt.Sprintf("SRID=%d;POINT(%v %v)", p.SRID, p.Lng, p.Lat)
}

// Scan :
func (p *Coord) Scan(val interface{}) error {
	if val == nil {
		return nil
	}

	var hexStr, ok = val.(string)
	if !ok {
		return errors.New("invalid input for boudary polygon")
	}

	rawHex, err := hex.DecodeString(hexStr)
	if nil != err {
		return err
	}

	g, err := ewkb.Unmarshal(rawHex)
	if nil != err {
		return err
	}
	var coords = g.FlatCoords()
	if len(coords) != 2 {
		return errors.New("invalid coord. ")
	}

	p.SRID = g.SRID()
	p.Lng = coords[0]
	p.Lat = coords[1]
	return nil
}

// Value :
func (p Coord) Value() (driver.Value, error) {
	return p.String(), nil
}

// MakeCoord :
func (p *Coord) MakeCoord() string {
	return fmt.Sprintf("ST_SetSRID(ST_MakeCoord(%f, %f), %d)", p.Lng, p.Lat, p.SRID)
}

func (p *Coord) GetCoord() geom.Coord {
	return geom.Coord{p.Lng, p.Lat}
}

func (p *Coord) To3857() *Coord {
	if p.SRID == 3857 {
		return NewCoord3857(p.Lng, p.Lat)
	}

	var p2 = NewCoord3857(
		(p.Lng*20037508.34)/180.0,
		math.Log(math.Tan(((90.0+p.Lat)*math.Pi)/360.0))/(math.Pi/180.0),
	)

	p2.Lat = (p2.Lat * 20037508.34) / 180
	return p2
}

func (p *Coord) To4326() *Coord {
	if p.SRID == 4326 {
		return NewCoord4326(p.Lng, p.Lat)
	}

	var p2 = NewCoord4326(
		(p.Lng*180.0)/20037508.34,
		(p.Lat*180.0)/20037508.34,
	)

	p2.Lat = (math.Atan(math.Pow(math.E, p2.Lat*(math.Pi/180.0)))*360.0)/math.Pi - 90.0
	return p2
}
