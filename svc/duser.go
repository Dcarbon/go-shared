package svc

import (
	"context"
	"github.com/Dcarbon/arch-proto/pb"
	"github.com/Dcarbon/go-shared/gutils"
	"time"
)

type IUserService interface {
	Init(ctx context.Context, request *pb.RUserInit) (*pb.RsUserInit, error)
	GetById(ctx context.Context, request *pb.RUserGetById) (*pb.UserInfo, error)
}
type userService struct {
	client pb.UserInfoServiceClient
}

func NewRequestUserService(host string) (IUserService, error) {
	cc, err := gutils.GetCCTimeout(host, 5*time.Second)
	if nil != err {
		return nil, err
	}

	var client = &userService{
		client: pb.NewUserInfoServiceClient(cc),
	}
	return client, nil
}

func (u userService) GetById(ctx context.Context, request *pb.RUserGetById) (*pb.UserInfo, error) {
	user, err := u.client.GetUserById(ctx, request)
	if err != nil {
		return nil, err
	}
	return &pb.UserInfo{
		Id:      user.Id,
		Phone:   user.Phone,
		Name:    user.Name,
		Email:   user.Email,
		Avatar:  user.Avatar,
		Address: user.Address,
		Gender:  user.Gender,
		Status:  user.Status,
		Dob:     user.Dob,
	}, nil
}

func (u userService) Init(ctx context.Context, request *pb.RUserInit) (*pb.RsUserInit, error) {
	user, err := u.client.Init(ctx, request)
	if err != nil {
		return nil, err
	}
	return &pb.RsUserInit{
		Id: user.Id,
	}, nil
}
