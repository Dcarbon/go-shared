package dmodels

import (
	"database/sql/driver"
	"fmt"
	"math/big"
)

type BigNumber struct {
	*big.Int
}

func NewBigNumber(x int64) *BigNumber {
	var b = &BigNumber{
		Int: big.NewInt(x),
	}
	return b
}

func NewBigNumberFromHex(text string) (*BigNumber, error) {
	var b = &BigNumber{
		Int: big.NewInt(0),
	}

	if text == "" {
		return b, nil
	}

	if len(text) > 2 && text[:2] == "0x" {
		n := big.NewInt(0)
		_, ok := n.SetString(text[2:], 16)
		if !ok {
			return nil, fmt.Errorf("invalid hex (%s) for bignumber", text)
		}
	}

	err := b.UnmarshalText([]byte(text))
	if nil != err {
		return nil, err
	}
	return b, nil
}

func MustNewBigNumberFromHex(text string) *BigNumber {
	var b, err = NewBigNumberFromHex(text)
	if nil != err {
		panic(err)
	}
	return b
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (b *BigNumber) UnmarshalJSON(val []byte) error {
	// Ignore null, like in the main JSON package.
	if string(val) == "null" {
		return nil
	}

	var text = string(val)
	if len(text) > 2 && text[:2] == "0x" {
		n := big.NewInt(0)
		_, ok := n.SetString(text[2:], 16)
		if !ok {
			return fmt.Errorf("invalid hex (%s) for bignumber", text)
		}
	}

	return b.UnmarshalText(val)
}

// MarshalJSON implements the json.Marshaler interface.
func (b *BigNumber) MarshalJSON() ([]byte, error) {
	if b == nil || b.Int == nil {
		return []byte("null"), nil
	}

	return []byte(b.ToHex()), nil
}

func (b *BigNumber) Scan(val interface{}) error {
	if nil == val {
		return nil
	}
	var err error
	var n = big.NewInt(0)
	switch val2 := val.(type) {
	case []byte:
		n.SetBytes(val2)
	default:
		err = fmt.Errorf("invalid type(%T) for bignumer", val2)
	}
	if nil != err {
		return err
	}

	if b == nil {
		b = new(BigNumber)
	}
	b.Int = n

	return nil
}

func (b *BigNumber) Value() (driver.Value, error) {
	if nil == b || b.Int == nil {
		return nil, nil
	}
	return b.Bytes(), nil
}

func (b *BigNumber) ToHex() string {
	var rs = fmt.Sprintf("%x", b.Int)
	if len(rs)%2 != 0 {
		return "0x0" + rs
	}
	return "0x" + rs
}
