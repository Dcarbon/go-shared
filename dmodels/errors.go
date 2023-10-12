package dmodels

import (
	"fmt"
	"log"
	"strings"

	"github.com/Dcarbon/go-shared/ecodes"
	"gorm.io/gorm"
)

var (
	ErrorUnauthorized     = NewError(ecodes.Unauthorized, "")
	ErrorPermissionDenied = NewError(ecodes.PermissionDenied, "")
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
} //@name Error

func NewError(code int, msg string) error {
	var err = &Error{
		Code:    code,
		Message: msg,
	}
	return err
}

func (err *Error) Error() string {
	return fmt.Sprintf("Code:%d Message:%s", err.Code, err.Message)
}

func (err Error) String() string {
	return fmt.Sprintf("Code:%d Message:%s", err.Code, err.Message)
}

func ParsePostgresError(label string, err error) error {
	if nil == err {
		return nil
	}
	log.Println("Postgres error: ", err)
	if err == gorm.ErrRecordNotFound {
		return NewError(
			ecodes.NotExisted,
			label+" is not existed",
		)
	}

	if strings.Contains(err.Error(), "duplicate") {
		return NewError(
			ecodes.Existed,
			label+" is existed",
		)
	}
	return ErrInternal(err)
}

// ErrInternal :
func ErrInternal(err error) error {
	if nil == err {
		return nil
	}
	log.Println("Internal error: ", err)
	return NewError(ecodes.Internal, "internal error")
}

func ErrNotFound(msg string) error {
	return NewError(ecodes.NotExisted, msg)
}

func ErrExisted(msg string, params ...interface{}) error {
	return NewError(ecodes.Existed, fmt.Sprintf(msg+" is existed", params...))
}

// ErrInternal :
func ErrNotImplement() error {
	return NewError(ecodes.NotImplement, "not implement")
}

// ErrInternal :
func ErrBadRequest(msg string) error {
	// log.Println("Bad request error: ", err)
	return NewError(ecodes.BadRequest, msg)
}

// ErrInternal :
func ErrQueryParam(msg string) error {
	return NewError(ecodes.QueryParamInvalid, msg)
}

func ErrInvalidSensorMetric(msg string) error {
	return NewError(ecodes.QueryParamInvalid, msg)
}

func ErrInvalidSignature() error {
	return NewError(ecodes.InvalidSignature, "Data of signature must be hex")
}

func ErrInvalidNonce() error {
	return NewError(ecodes.IOTInvalidNonce, "Invalid nonce")
}
