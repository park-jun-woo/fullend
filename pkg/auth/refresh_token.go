package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// @func refreshToken
// @description 리프레시 토큰을 발급한다 (7일 만료)

type RefreshTokenInput struct {
	UserID int64
	Email  string
	Role   string
}

type RefreshTokenOutput struct {
	RefreshToken string
}

func RefreshToken(in RefreshTokenInput) (RefreshTokenOutput, error) {
	claims := jwt.MapClaims{
		"user_id": in.UserID,
		"email":   in.Email,
		"role":    in.Role,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte("secret"))
	return RefreshTokenOutput{RefreshToken: signed}, err
}
