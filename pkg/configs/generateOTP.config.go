package pkg

import (
	"math/rand"
	"time"
)

func GenerateOTP() string {

	rand.Seed(time.Now().UnixNano())

	otp := string(rune(rand.Intn(9000) + 1000))

	return otp

}
