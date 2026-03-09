package crypto

import "github.com/pquerna/otp/totp"

// @func generateOTP
// @description TOTP 시크릿과 QR 프로비저닝 URL을 생성한다

type GenerateOTPInput struct {
	Issuer      string
	AccountName string
}

type GenerateOTPOutput struct {
	Secret string
	URL    string // otpauth:// URL (QR 코드용)
}

func GenerateOTP(in GenerateOTPInput) (GenerateOTPOutput, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      in.Issuer,
		AccountName: in.AccountName,
	})
	if err != nil {
		return GenerateOTPOutput{}, err
	}
	return GenerateOTPOutput{
		Secret: key.Secret(),
		URL:    key.URL(),
	}, nil
}
