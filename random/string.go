package random

import (
	"math/rand"
	"time"
)

const randCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandString(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = randCharset[r.Intn(len(randCharset))]
	}

	return string(b)
}
