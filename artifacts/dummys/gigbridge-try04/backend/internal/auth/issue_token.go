//fullend:gen ssot=fullend.yaml contract=811d966
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
	ID int64
	Role string
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
		"user_id": req.ID,
		"role": req.Role,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	return IssueTokenResponse{AccessToken: signed}, err
}
