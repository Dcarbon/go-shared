package gutils

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var ctxKey = new(int)

// Check authen and permission
type A2M struct {
	jwtKey string
	perm   string
}

func NewA2(jwtKey string, perm string) *A2M {
	var a2 = &A2M{
		jwtKey: jwtKey,
		perm:   perm,
	}
	return a2
}

func (a2 *A2M) HandlerFunc(r *gin.Context) {
	var authToken = r.GetHeader("Authorization")
	var idx = strings.Index(authToken, "Bearer ")
	if idx != 0 && len(authToken) < 10 {
		r.AbortWithError(http.StatusUnauthorized, ErrorUnauthorized)
		return
	}

	var user, err = DecodeJWT(a2.jwtKey, authToken[7:])
	if nil != err {
		fmt.Println("Decode jwt error: ", err)
		r.AbortWithError(http.StatusUnauthorized, ErrorUnauthorized)
		return
	}

	// err = hasPerm(user.Role, a2.perm)
	// if nil != err {
	// 	r.AbortWithError(http.StatusUnauthorized, ErrorUnauthorized)
	// 	return
	// }

	var ctx = context.WithValue(r.Request.Context(), ctxKey, user)
	r.Request = r.Request.WithContext(ctx)
}
