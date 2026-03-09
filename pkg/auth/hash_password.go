package auth

import "golang.org/x/crypto/bcrypt"

// @func hashPassword
// @description 평문 비밀번호를 bcrypt 해시로 변환한다

type HashPasswordInput struct {
	Password string
}

type HashPasswordOutput struct {
	HashedPassword string
}

func HashPassword(in HashPasswordInput) (HashPasswordOutput, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	return HashPasswordOutput{HashedPassword: string(hash)}, err
}
