package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// @func issueToken
// @description 인증된 사용자 정보로 JWT 액세스 토큰을 발급한다

type IssueTokenInput struct {
	UserID int64
	Email  string
	Role   string
}

type IssueTokenOutput struct {
	AccessToken string
}

func IssueToken(in IssueTokenInput) (IssueTokenOutput, error) {
	claims := jwt.MapClaims{
		"user_id": in.UserID,
		"email":   in.Email,
		"role":    in.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte("secret"))
	return IssueTokenOutput{AccessToken: signed}, err
}
