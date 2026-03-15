//fullend:gen ssot=fullend.yaml contract=78bda02
package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// @func refreshToken
// @description 리프레시 토큰을 발급한다 (7일 만료)

type RefreshTokenRequest struct {
	Email string
	Role string
	UserID int64
}

type RefreshTokenResponse struct {
	RefreshToken string
}

func RefreshToken(req RefreshTokenRequest) (RefreshTokenResponse, error) {
	claims := jwt.MapClaims{
		"email": req.Email,
		"role": req.Role,
		"user_id": req.UserID,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte("secret"))
	return RefreshTokenResponse{RefreshToken: signed}, err
}
