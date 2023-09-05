package gutils

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func GetCC(host string) (*grpc.ClientConn, error) {
	if host == "" {
		return nil, ErrBadRequest("Can't connect with host is empty")
	}
	cc, err := grpc.Dial(
		host,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if nil != err {
		return nil, err
	}
	return cc, err
}

func MustGetCC(host string) *grpc.ClientConn {
	if host == "" {
		panic(ErrBadRequest("Can't connect with host is empty"))
	}

	cc, err := grpc.Dial(
		host,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if nil != err {
		log.Fatalf("Dial to %s error: %s\n", host, err.Error())
	}
	return cc
}

func GetCCTimeout(host string, dur time.Duration) (*grpc.ClientConn, error) {
	if host == "" {
		return nil, ErrBadRequest("Can't connect with host is empty")
	}

	var ctx, cancel = context.WithTimeout(context.TODO(), dur)
	var cc, err = grpc.DialContext(
		ctx,
		host,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	defer cancel()
	if nil != err {
		return nil, err
	}
	return cc, ctx.Err()
}

// MetaWithJWT :
func MetaWithJWT(ctx context.Context, jwt string) context.Context {
	var md = metadata.New(
		map[string]string{
			"Authorization": "Bearer " + jwt,
		},
	)
	return metadata.NewOutgoingContext(ctx, md)
}
