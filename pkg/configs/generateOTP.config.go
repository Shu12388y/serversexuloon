package pkg

import (
	"math/rand"
	"time"
)

// GenerateOTP returns a 4-digit OTP as a string
func GenerateOTP() int {
	rand.Seed(time.Now().UnixNano())
	otp := rand.Intn(9000) + 1000
	return otp
}
