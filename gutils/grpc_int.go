package gutils

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
)

func catchError(info *grpc.UnaryServerInfo) {
	err := recover()
	if nil != err {
		fmt.Printf(
			"Criticial crash path:%s error:%+v \n",
			info.FullMethod,
			err,
		)
	}
}

// UnaryPreventPanic :
func UnaryPreventPanic(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	defer catchError(info)
	return handler(ctx, req)
}

// ARConfig : auth request config
type ARConfig struct {
	Require     bool
	Path        string
	Permission  string
	PermDesc    string
	DefaultRole []string
}

// LogInterceptor :
type LogInterceptor struct {
}

// NewLogInterceptor :
func NewLogInterceptor() *LogInterceptor {
	return &LogInterceptor{}
}

// Intercept :
func (ai *LogInterceptor) Intercept(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	data, err := handler(ctx, req)
	if nil != err {
		log.Println(info.FullMethod+" error: ", err)
	} else {
		log.Println(info.FullMethod + " success")
	}
	return data, err
}
