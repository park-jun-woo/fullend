//ff:func feature=gen-gogin type=generator control=iteration dimension=1
//ff:what creates internal/auth/refresh_token.go with 7-day refresh token generation

package gogin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/fullend/pkg/contract"
	"github.com/park-jun-woo/fullend/pkg/parser/manifest"
)

// generateRefreshToken creates internal/auth/refresh_token.go.
func generateRefreshToken(authDir string, claims map[string]manifest.ClaimDef, fields []string, claimsHash string) error {
	var reqFields []string
	var mapClaims []string
	for _, field := range fields {
		def := claims[field]
		reqFields = append(reqFields, fmt.Sprintf("\t%s %s", field, def.GoType))
		mapClaims = append(mapClaims, fmt.Sprintf("\t\t\"%s\": req.%s,", def.Key, field))
	}

	src := fmt.Sprintf(`package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// @func refreshToken
// @description 리프레시 토큰을 발급한다 (7일 만료)

type RefreshTokenRequest struct {
%s
}

type RefreshTokenResponse struct {
	RefreshToken string
}

func RefreshToken(req RefreshTokenRequest) (RefreshTokenResponse, error) {
	claims := jwt.MapClaims{
%s
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte("secret"))
	return RefreshTokenResponse{RefreshToken: signed}, err
}
`, strings.Join(reqFields, "\n"), strings.Join(mapClaims, "\n"))

	d := &contract.Directive{Ownership: "gen", SSOT: "fullend.yaml", Contract: claimsHash}
	src = injectFileDirective(src, d)
	return os.WriteFile(filepath.Join(authDir, "refresh_token.go"), []byte(src), 0644)
}
