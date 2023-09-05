package gutils

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// AuthInterceptor :
type AuthInterceptor struct {
	jwtKey string
	config map[string]*ARConfig
	CtxKey *int
	// iam    adin.IAMClient
}

// NewAuthInterceptor :
func NewAuthInterceptor(
	iamHost, jwtKey string,
	config map[string]*ARConfig,
) (*AuthInterceptor, error) {
	// var iam adin.IAMClient
	// if iamHost != "" {
	// 	iam = GetIAMBlock(iamHost)
	// }

	var ai = &AuthInterceptor{
		// iam:    iam,
		jwtKey: jwtKey,
		config: config,
		CtxKey: new(int),
	}
	return ai, nil
}

func (ai *AuthInterceptor) GetContextKey() interface{} {
	return ai.CtxKey
}

// Intercept :
func (ai *AuthInterceptor) Intercept(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	var config = ai.config[info.FullMethod]
	if nil == config {
		fmt.Printf("Method %s has no config \n", info.FullMethod)
		return nil, ErrorNotRegisterAuth
	}
	if config.Require {
		staff, err := ai.authen(ctx)
		if nil != err {
			return nil, ErrorUnauthorized
		}
		err = ai.checkPerm(ctx, staff, config.Permission)
		if nil != err {
			return nil, ErrorNoPermission
		}
		ctx = context.WithValue(ctx, ai.CtxKey, staff)
	}
	return handler(ctx, req)
}

func (ai *AuthInterceptor) checkPerm(ctx context.Context, staff *ClaimModel, permID string,
) error {
	// has no iam, only for iam service
	// if nil == ai.iam || permID == "" {
	// 	return nil
	// }
	// if staff.Role == "isv" || staff.Role == "super-admin" {
	// 	return nil
	// }
	// resp, err := ai.iam.Permit(
	// 	context.TODO(),
	// 	&adin.PermitRequest{
	// 		RoleID: staff.Role,
	// 		PermID: permID,
	// 	},
	// )
	// if nil != err {
	// 	return err
	// }
	// if !resp.Permit {
	// 	log.Printf("User %d try to execute perm %s \n", staff.Id, permID)
	// 	return ErrorNoPermission
	// }
	return nil
}

func (ai *AuthInterceptor) authen(ctx context.Context) (*ClaimModel, error) {
	var headers, ok = metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, ErrBadRequest("Cant get headers")
	}

	var authToken = headers.Get("Authorization")
	if len(authToken) == 0 {
		return nil, ErrorUnauthorized
	}
	var idx = strings.Index(authToken[0], "Bearer ")
	if idx != 0 {
		return nil, ErrorUnauthorized
	}

	var staff, err = DecodeJWT(ai.jwtKey, authToken[0][7:])
	if nil != err {
		fmt.Println("Decode jwt error: ", err)
		return nil, ErrorUnauthorized
	}
	return staff, nil
}

func (ai *AuthInterceptor) GetAuth(ctx context.Context) (*ClaimModel, error) {
	staff, ok := ctx.Value(ai.CtxKey).(*ClaimModel)
	if !ok {
		return nil, ErrorNotRegisterAuth
	}
	return staff, nil
}
