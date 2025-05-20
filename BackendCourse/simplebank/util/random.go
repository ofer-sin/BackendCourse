package util

import (
	"math/rand"
	"time"
)

const alpghabet = "abcdefghijklmnopqrstuvwxyz"

var randomGen = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomInt(min, max int64) int64 {
	if min >= max {
		panic("invalid min and max")
	}
	return min + randomGen.Int63n(max-min+1)
}

func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alpghabet[randomGen.Intn(len(alpghabet))]
	}
	return string(b)
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD"}
	n := len(currencies)
	return currencies[randomGen.Intn(n)]
}
