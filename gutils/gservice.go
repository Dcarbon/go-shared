package gutils

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/metadata"
)

type GService struct {
	internalJWT string
	Config      Config
}

func NewGService(config Config, jwt string) (*GService, error) {
	var gs = &GService{
		Config: config,
	}
	var err error
	if jwt == "" {
		err = gs.login()
	} else {
		gs.internalJWT = jwt
	}

	return gs, err
}

func (gs *GService) login() error {
	var user = gs.Config.Options[ISVUser]
	var claim = &customClaim{
		Auth: &ClaimModel{
			Id:        2,
			Role:      "isv",
			FirstName: "",
			LastName:  "",
			Username:  user,
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: int64(10*365*86400) + int64(time.Now().Unix()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	jwt, err := token.SignedString([]byte(gs.Config.JwtKey))
	if nil != err {
		return err
	}

	if nil != err {
		return err
	}
	gs.internalJWT = jwt
	return nil
}

// func (gs *GService) CreateConfigPerm() error {
// 	var cc, err = GetCCTimeout(gs.Config.GetIAM(), 10*time.Second)
// 	if nil != err {
// 		return err
// 	}
// 	defer cc.Close()
// 	var iamClient = adin.NewIAMClient(cc)
// 	var ctx = gs.WithJwt(context.TODO())
// 	for path, authConfig := range gs.Config.AuthConfig {
// 		if !authConfig.Require || authConfig.Permission == "" {
// 			continue
// 		}
// 		if authConfig.PermDesc == "" {
// 			return fmt.Errorf("path:%s missing perm desc", path)
// 		}
// 		_, err = iamClient.PermCreate(ctx, &adin.PermModel{
// 			Id:    authConfig.Permission,
// 			Group: gs.Config.Name,
// 			Desc:  authConfig.PermDesc,
// 		})
// 		if nil != err {
// 			return err
// 		}
// 		_, err = iamClient.Assign(
// 			ctx,
// 			&adin.AssignRequest{
// 				RoleID: "super-admin",
// 				PermID: authConfig.Permission,
// 			})
// 		if nil != err {
// 			return fmt.Errorf(
// 				"assign perm:%s for super-admin error: %s",
// 				authConfig.Permission, err.Error(),
// 			)
// 		}
// 		for _, r := range authConfig.DefaultRole {
// 			_, err = iamClient.Assign(
// 				ctx,
// 				&adin.AssignRequest{
// 					RoleID: r,
// 					PermID: authConfig.Permission,
// 				})
// 			if nil != err {
// 				return fmt.Errorf(
// 					"assign perm:%s for role:%s error: %s",
// 					authConfig.Permission, r, err.Error(),
// 				)
// 			}
// 		}
// 	}
// 	return nil
// }

func (gs *GService) GetOption(key string) string {
	return gs.Config.Options[key]
}

func (gs *GService) GetToken() string {
	return gs.internalJWT
}

func (gs *GService) WithJwt(ctx context.Context) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.New(
		map[string]string{
			"Authorization": "Bearer " + gs.internalJWT,
		},
	))
}
