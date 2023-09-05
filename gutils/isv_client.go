package gutils

import (
	"context"
	"log"

	"github.com/Dcarbon/go-shared/libs/utils"
	"google.golang.org/grpc/metadata"
)

type ISVClient struct {
	internalJWT string //
	name        string // Client name
	user        string
	pwd         string
	authHost    string
}

func NewInternalClient(name, authHost, isvUser, isvPwd string,
) (*ISVClient, error) {
	var cc = &ISVClient{
		name:     name,
		user:     isvUser,
		pwd:      isvPwd,
		authHost: authHost,
	}
	return cc, cc.Login()
}

func NewInternalClientFromEnv(name string) (*ISVClient, error) {
	var cc = &ISVClient{
		name:     name,
		user:     utils.StringEnv(ISVUser, "isv"),
		pwd:      utils.StringEnv(ISVPass, "a1b2cba"),
		authHost: utils.StringEnv(ISVKeyStaff, "10.60.0.50:9005"),
	}
	return cc, cc.Login()
}

func NewInternalClientFromToken(name, token string) *ISVClient {
	if token == "" {
		panic("Invalid internal token")
	}
	return &ISVClient{
		internalJWT: token,
		name:        name,
	}
}

func (ic *ISVClient) Login() error {
	// var user = ic.user
	// var pass = ic.pwd

	// c3, err := GetCCTimeout(ic.authHost, 5*time.Second)
	// if nil != err {
	// 	return err
	// }
	// defer c3.Close()

	// data, err := adin.NewStaffClient(c3).Login(
	// 	context.TODO(),
	// 	&adin.RLogin{
	// 		Username: user,
	// 		Password: pass,
	// 	})
	// if nil != err {
	// 	return err
	// }
	// ic.internalJWT = data.Token
	return nil
}

func (ic *ISVClient) GetJWT() string {
	return ic.internalJWT
}

func (ic *ISVClient) WithJwt(ctx context.Context) context.Context {
	if ic.internalJWT == "" {
		err := ic.Login()
		if nil != err {
			log.Printf("%s login error: %s\n", ic.name, err.Error())
		}
	}

	return metadata.NewOutgoingContext(ctx, metadata.New(
		map[string]string{
			"Authorization": "Bearer " + ic.internalJWT,
		},
	))
}
