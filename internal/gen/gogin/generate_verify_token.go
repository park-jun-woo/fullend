//ff:func feature=gen-gogin type=generator control=iteration dimension=1
//ff:what creates internal/auth/verify_token.go with JWT verification and claims extraction

package gogin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/fullend/internal/contract"
	"github.com/park-jun-woo/fullend/internal/projectconfig"
)

// generateVerifyToken creates internal/auth/verify_token.go.
func generateVerifyToken(authDir string, claims map[string]projectconfig.ClaimDef, fields []string, claimsHash string) error {
	var respFields []string
	var extractLines []string
	for _, field := range fields {
		def := claims[field]
		respFields = append(respFields, fmt.Sprintf("\t%s %s", field, def.GoType))
		extractLines = append(extractLines, claimExtractLine(field, def))
	}

	var resultFields []string
	for _, field := range fields {
		def := claims[field]
		resultFields = append(resultFields, resultAssignLine(field, def))
	}

	src := fmt.Sprintf(`package auth

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
%s
}

func VerifyToken(req VerifyTokenRequest) (VerifyTokenResponse, error) {
	token, err := jwt.Parse(req.Token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %%v", t.Header["alg"])
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
%s
	return VerifyTokenResponse{
%s
	}, nil
}
`, strings.Join(respFields, "\n"), strings.Join(extractLines, "\n"), strings.Join(resultFields, "\n"))

	d := &contract.Directive{Ownership: "gen", SSOT: "fullend.yaml", Contract: claimsHash}
	src = injectFileDirective(src, d)
	return os.WriteFile(filepath.Join(authDir, "verify_token.go"), []byte(src), 0644)
}
