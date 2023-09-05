package gutils

import (
	"fmt"
	"testing"

	"github.com/Dcarbon/go-shared/libs/utils"
)

func TestJWT(t *testing.T) {
	var key = "1ldsfjsldfjsldfjlj3l1k3j1l2jp09gsdvsldkvnslvnsldnlnlnjsdflsjdf"
	var model = &ClaimModel{
		Id:        100,
		Role:      "ronlo",
		FirstName: "dips",
		LastName:  "urgf",
		Username:  "user_test",
	}

	var token, err = EncodeJWTClaim(key, model)
	fmt.Println(token)
	utils.PanicError("Encode jwt error", err)

	encodeModel, err := DecodeJWT(key, token)
	utils.PanicError("Decode jwt error", err)
	fmt.Println("Encode model: ", encodeModel)
}
