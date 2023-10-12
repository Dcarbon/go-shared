package gutils

import (
	"fmt"
	"log"
	"strings"

	"github.com/Dcarbon/go-shared/ecodes"
	"github.com/jackc/pgconn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// Define exist
var (
	ErrorUnauthorized = status.Errorf(ecodes.Unauthorized, "unauthorized")
	ErrorNoPermission = status.Errorf(ecodes.PermissionDenied, "has no permission")

	ErrorNotImplement    = status.Errorf(ecodes.NotImplement, "not implement")
	ErrorNotRegisterAuth = status.Errorf(ecodes.NotRegisterAuth, "not register auth")

	// ErrorStaffIsDisable = status.Errorf(ecodes.StaffIsDisable, "was disabled")
)

func NewError(code int, msg string) error {
	return status.Errorf(codes.Code(code), msg)
}

// ErrInternal :
func ErrInternal(err error) error {
	if nil == err {
		return nil
	}
	log.Println("Internal error: ", err)
	return status.Errorf(ecodes.Internal, "internal error")
}

// ErrBadRequest :
func ErrBadRequest(msg string) error {
	return status.Errorf(ecodes.BadRequest, msg)
}

// ErrBadRequest :
func ErrBadRequestf(msg string, params ...interface{}) error {
	return status.Errorf(ecodes.BadRequest, fmt.Sprintf(msg, params...))
}

// ErrNotFound :
func ErrNotFound(msg string) error {
	return status.Errorf(ecodes.NotExisted, msg)
}

// ErrNotFound :
func ErrNotFoundf(msg string, params ...interface{}) error {
	return status.Errorf(ecodes.NotExisted, fmt.Sprintf(msg, params...))
}

func ErrExisted(msg string, params ...interface{}) error {
	return status.Errorf(ecodes.Existed, fmt.Sprintf(msg, params...))
}

// func ParsePostgresError(err error, modelName string) error {
// 	if nil == err {
// 		return nil
// 	}

// 	if err == gorm.ErrRecordNotFound {
// 		return status.Errorf(
// 			ecodes.NotExisted,
// 			modelName+" is not existed",
// 		)
// 	}

// 	if strings.Contains(err.Error(), "duplicate") {
// 		return status.Errorf(
// 			ecodes.NotExisted,
// 			modelName+" is existed",
// 		)
// 	}
// 	return ErrInternal(err)
// }

func ParsePostgres(modelName string, err error) error {
	if nil == err {
		return nil
	}

	if err == gorm.ErrRecordNotFound {
		return status.Errorf(
			ecodes.NotExisted,
			modelName+" is not existed",
		)
	}

	pgErr, ok := err.(*pgconn.PgError)
	if ok {
		return status.Errorf(
			ecodes.NotExisted,
			modelName+pgErr.Detail,
		)
	}

	if strings.Contains(err.Error(), "duplicate") {
		return status.Errorf(
			ecodes.NotExisted,
			modelName+" is existed",
		)
	}
	return ErrInternal(err)
}

// IsPostgresNotFound :
func IsPostgresNotFound(err error) error {
	if nil == err {
		return nil
	}
	if err == gorm.ErrRecordNotFound {
		return status.Errorf(ecodes.NotExisted, "model is not existed")
	}
	return ErrInternal(err)
}

func IsPostgresDuplicate(err error) error {
	if nil == err {
		return nil
	}
	if strings.Contains(err.Error(), "duplicate") {
		return status.Errorf(ecodes.NotExisted, "model is existed")
	}
	return ErrInternal(err)
}
