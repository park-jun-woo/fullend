package auth

import (
	"crypto/rand"
	"encoding/hex"
)

// @func generateResetToken
// @description 비밀번호 리셋용 32바이트 랜덤 hex 토큰을 생성한다

type GenerateResetTokenInput struct{}

type GenerateResetTokenOutput struct {
	Token string
}

func GenerateResetToken(in GenerateResetTokenInput) (GenerateResetTokenOutput, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return GenerateResetTokenOutput{}, err
	}
	return GenerateResetTokenOutput{Token: hex.EncodeToString(b)}, nil
}
