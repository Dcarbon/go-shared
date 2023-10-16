package gutils

import (
	"encoding/json"
	"fmt"
)

func Dump(data interface{}) {
	a, err := json.MarshalIndent(data, "", "")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(a))
}
