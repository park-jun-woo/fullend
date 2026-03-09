package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// @func verifyToken
// @description JWT 토큰을 검증하고 claims를 추출한다

type VerifyTokenInput struct {
	Token  string
	Secret string
}

type VerifyTokenOutput struct {
	UserID int64
	Email  string
	Role   string
}

func VerifyToken(in VerifyTokenInput) (VerifyTokenOutput, error) {
	token, err := jwt.Parse(in.Token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(in.Secret), nil
	})
	if err != nil {
		return VerifyTokenOutput{}, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return VerifyTokenOutput{}, fmt.Errorf("invalid token")
	}
	userID, _ := claims["user_id"].(float64)
	email, _ := claims["email"].(string)
	role, _ := claims["role"].(string)
	return VerifyTokenOutput{
		UserID: int64(userID),
		Email:  email,
		Role:   role,
	}, nil
}
