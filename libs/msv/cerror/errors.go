package cerror

import (
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrorUnauthorized     = NewError(ECodeUnauthorized, "unauthorized")
	ErrorPermissionDenied = NewError(ECodePermissionDenied, "permission denied")
)

type Error struct {
	Code    ECode  `json:"code"`
	Message string `json:"message"`
}

func NewError(code ECode, msg string) error {
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
			ECodeNotExisted,
			label+" is not existed",
		)
	}

	if strings.Contains(err.Error(), "duplicate") {
		return NewError(
			ECodeExisted,
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
	return NewError(ECodeInternal, "internal error")
}

func ErrNotImplement() error {
	return NewError(ECodeNotImplement, "not implement")
}

func ErrBadRequest(msg string) error {
	// log.Println("Bad request error: ", err)
	return NewError(ECodeBadRequest, msg)
}

func ErrQueryParam(msg string) error {
	return NewError(ECodeQueryParamInvalid, msg)
}
