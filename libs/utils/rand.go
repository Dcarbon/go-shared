package utils

import (
	"math/rand"
	"time"
)

const letterRunes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

func init() {
	rand.Seed(time.Now().Unix())
}

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
