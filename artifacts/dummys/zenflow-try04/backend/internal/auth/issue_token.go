//fullend:gen ssot=fullend.yaml contract=78bda02
package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// @func issueToken
// @description claims 기반 JWT 액세스 토큰을 발급한다

type IssueTokenRequest struct {
	Email string
	Role string
	UserID int64
}

type IssueTokenResponse struct {
	AccessToken string
}

func IssueToken(req IssueTokenRequest) (IssueTokenResponse, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret"
	}
	claims := jwt.MapClaims{
		"email": req.Email,
		"role": req.Role,
		"user_id": req.UserID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	return IssueTokenResponse{AccessToken: signed}, err
}
