package model

// Token은 인증 토큰이다.
type Token struct {
	AccessToken string
	ExpiresAt   string
}

// issueToken은 JWT 토큰을 발급한다.
// SSaC에서 @func issueToken으로 참조된다.
func issueToken(userID int64) (Token, error) { return Token{}, nil }

// hashPassword는 비밀번호를 해싱한다.
// SSaC에서 @func hashPassword로 참조된다.
func hashPassword(plain string) (string, error) { return "", nil }
