package utils

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghmnopqrstuvwxyz"

func init() {
	rand.New(rand.NewSource(time.Now().UnixMicro()))
}

// RandomInt value in [min, max]
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString
func RandomString(n int) string {
	var sb strings.Builder

	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(0, 10000)
}

func RandomCurrency() string {
	currencies := []string{"USD", "VND", "EUR"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
