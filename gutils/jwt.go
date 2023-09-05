package gutils

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type ClaimModel struct {
	Id        int64
	Role      string
	FirstName string
	LastName  string
	Username  string
}

type customClaim struct {
	jwt.StandardClaims
	Auth *ClaimModel
}

// DecodeJWT :
func DecodeJWT(key string, token string) (*ClaimModel, error) {
	var claim = &customClaim{}
	jtoken, err := jwt.ParseWithClaims(
		token,
		claim,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		},
	)
	if nil != err {
		return nil, err
	}

	if !jtoken.Valid {
		return nil, errors.New("token is invalid")
	}
	return claim.Auth, nil
}

// EncodeJWT :
func EncodeJWT(key string, model *ClaimModel) (string, error) {
	var claim = &customClaim{
		Auth: &ClaimModel{
			Id:        model.Id,
			Role:      model.Role,
			FirstName: model.FirstName,
			LastName:  model.LastName,
			Username:  model.Username,
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 86400,
		},
	}
	var token = jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(key))
}

// EncodeJWT :
func EncodeJWTClaim(key string, model *ClaimModel) (string, error) {
	var claim = &customClaim{
		Auth: model,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 86400,
		},
	}
	var token = jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(key))
}

// EncodeJWTWithExpire :
func EncodeJWTWithExpire(key string, model *ClaimModel, expire int64) (string, error) {
	var claim = &customClaim{
		Auth: &ClaimModel{
			Id:        model.Id,
			Role:      model.Role,
			FirstName: model.FirstName,
			LastName:  model.LastName,
			Username:  model.Username,
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire,
		},
	}
	var token = jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(key))
}

// DecodeJWTRequest :
func DecodeJWTRequest(key string, r *http.Request) (*ClaimModel, error) {
	var auth = r.Header.Get("Authorization")
	var idx = strings.Index(auth, "Bearer ")
	if idx != 0 {
		return nil, errors.New("Unauthorized")
	}

	var staff, err = DecodeJWT(key, auth[7:])
	if nil != err {
		return nil, errors.New("Unauthorized")
	}
	return staff, nil
}
