//fullend:gen ssot=fullend.yaml contract=78bda02
package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// @func verifyToken
// @error 401
// @description JWT 토큰을 검증하고 claims를 추출한다

type VerifyTokenRequest struct {
	Token  string
	Secret string
}

type VerifyTokenResponse struct {
	Email string
	Role string
	UserID int64
}

func VerifyToken(req VerifyTokenRequest) (VerifyTokenResponse, error) {
	token, err := jwt.Parse(req.Token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(req.Secret), nil
	})
	if err != nil {
		return VerifyTokenResponse{}, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return VerifyTokenResponse{}, fmt.Errorf("invalid token")
	}
	emailRaw, _ := claims["email"].(string)
	roleRaw, _ := claims["role"].(string)
	var userID int64
	if v, ok := claims["user_id"].(float64); ok {
		userID = int64(v)
	}
	return VerifyTokenResponse{
		Email: emailRaw,
		Role: roleRaw,
		UserID: userID,
	}, nil
}
