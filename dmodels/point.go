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

// Point:
type Point struct {
	SRID int     `json:"srid"` //
	Lat  float64 `json:"lat"`  // vi tuyen (pgis: y)
	Lng  float64 `json:"lng"`  // kinh tuyen:(pgis: x)
}

func NewPoint(srid int, lng, lat float64) *Point {
	return &Point{
		SRID: srid,
		Lat:  lat,
		Lng:  lng,
	}
}

func NewPoint4326(lng, lat float64) *Point {
	return &Point{
		SRID: 4326,
		Lat:  lat,
		Lng:  lng,
	}
}

func NewPoint3857(lng, lat float64) *Point {
	return &Point{
		SRID: 3857,
		Lat:  lat,
		Lng:  lng,
	}
}

func (p *Point) String() string {
	return fmt.Sprintf("SRID=%d;POINT(%v %v)", p.SRID, p.Lng, p.Lat)
}

// Scan :
func (p *Point) Scan(val interface{}) error {
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
func (p Point) Value() (driver.Value, error) {
	return p.String(), nil
}

// MakePoint :
func (p *Point) MakePoint() string {
	return fmt.Sprintf("ST_SetSRID(ST_MakePoint(%f, %f), %d)", p.Lng, p.Lat, p.SRID)
}

func (p *Point) GetCoord() geom.Coord {
	return geom.Coord{p.Lng, p.Lat}
}

func (p *Point) To3857() *Point {
	if p.SRID == 3857 {
		return NewPoint3857(p.Lng, p.Lat)
	}

	var p2 = NewPoint3857(
		(p.Lng*20037508.34)/180,
		math.Log(math.Tan(((90.0+p.Lat)*math.Pi)/360))/(math.Pi/180),
	)

	p2.Lat = (p2.Lat * 20037508.34) / 180
	return p2
}

func (p *Point) To4326() *Point {
	if p.SRID == 4326 {
		return NewPoint4326(p.Lng, p.Lat)
	}

	var p2 = NewPoint4326(
		(p.Lng*180.0)/20037508.34,
		(p.Lat*180.0)/20037508.34,
	)

	p2.Lat = (math.Atan(math.Pow(math.E, p2.Lat*(math.Pi/180.0)))*360.0)/math.Pi - 90.0
	return p2
}
