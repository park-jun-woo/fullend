package crypto

import (
	"fmt"

	"github.com/pquerna/otp/totp"
)

// @func verifyOTP
// @description TOTP 코드가 시크릿과 일치하는지 검증한다

type VerifyOTPInput struct {
	Code   string
	Secret string
}

type VerifyOTPOutput struct{}

func VerifyOTP(in VerifyOTPInput) (VerifyOTPOutput, error) {
	if !totp.Validate(in.Code, in.Secret) {
		return VerifyOTPOutput{}, fmt.Errorf("invalid OTP code")
	}
	return VerifyOTPOutput{}, nil
}
