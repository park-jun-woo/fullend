//fullend:gen ssot=fullend.yaml contract=445e7aa
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
	ID int64
	OrgID int64
	Role string
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
	idRaw, _ := claims["user_id"].(float64)
	orgIDRaw, _ := claims["org_id"].(float64)
	roleRaw, _ := claims["role"].(string)
	return VerifyTokenResponse{
		Email: emailRaw,
		ID: int64(idRaw),
		OrgID: int64(orgIDRaw),
		Role: roleRaw,
	}, nil
}
