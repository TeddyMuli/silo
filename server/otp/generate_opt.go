package otp

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// GenerateOTP generates a 6-digit OTP.
func GenerateOTP() (string, error) {
	const otpLength = 6
	otp := ""

	for i := 0; i < otpLength; i++ {
		digit, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		otp += fmt.Sprintf("%d", digit)
	}

	return otp, nil
}
