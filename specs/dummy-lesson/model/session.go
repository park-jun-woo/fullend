package model

// @dto
// Token은 인증 토큰이다 (DDL 테이블 없음).
type Token struct {
	AccessToken string
	ExpiresAt   string
}
