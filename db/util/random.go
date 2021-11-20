package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generate a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max - min + 1)
}

// RandomString generate a random string of length n
func RandomString(n int) string{
	var sb strings.Builder
	k := len(alphabet)
	for i:=0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// RandomOwner generate a random owner
func RandomOwner() string {
	return RandomString(6)
}

// RandomBalance generate a random balance
func RandomBalance() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency generate a random currency
func RandomCurrency() string {
	currency := []string{"USD", "ERU", "RMB"}
	k := len(currency)
	return currency[rand.Intn(k)]
}

// RandomMoney generate random money value
func RandomMoney() int64 {
	return RandomInt(-1000, 1000)
}

// RandomAmount generaete positive integer
func RandomAmount() int64 {
	return RandomInt(0, 1000)
}

// RandomEmail generaete a legal random google email
func RandomEmail() string {
	emailname := RandomOwner()
	emailType := "@gmail.com"
	return emailname + emailType
}