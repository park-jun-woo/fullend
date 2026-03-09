package auth

import "golang.org/x/crypto/bcrypt"

// @func verifyPassword
// @description 저장된 해시와 평문 비밀번호가 일치하는지 검증한다

type VerifyPasswordInput struct {
	PasswordHash string
	Password     string
}

type VerifyPasswordOutput struct{}

func VerifyPassword(in VerifyPasswordInput) (VerifyPasswordOutput, error) {
	err := bcrypt.CompareHashAndPassword([]byte(in.PasswordHash), []byte(in.Password))
	return VerifyPasswordOutput{}, err
}
