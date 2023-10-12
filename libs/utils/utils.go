package utils

import (
	"encoding/json"
	"log"
)

// PanicError :
func PanicError(msg string, err error, data ...interface{}) {
	if nil != err {
		panic(msg + " error: " + err.Error())
	}

	if len(data) == 0 {
		return
	}

	for _, it := range data {
		if it == nil {
			log.Println(msg + " is null")
		} else {
			raw, _ := json.MarshalIndent(it, "", "  ")
			log.Println(msg, ":", string(raw))
		}
	}
}

// Dump : use for test
func Dump(label string, data interface{}) {
	if nil == data {
		log.Println(label+" dump: ", nil)
		return
	}
	var raw, err = json.MarshalIndent(data, "", "  ")
	if nil != err {
		panic(err)
	}
	log.Println(label + " dump: " + string(raw))
}
