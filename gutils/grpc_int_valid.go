package gutils

import (
	"context"
	"errors"

	"google.golang.org/grpc"
)

type ValidatorInterceptor struct {
}

//Intercept :
func (ai *ValidatorInterceptor) Intercept(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	return nil, ErrInternal(errors.New("validator is not implement"))
}
