//ff:func feature=gen-gogin type=generator
//ff:what creates internal/auth/issue_token.go with claims-based JWT issuing

package gogin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/geul-org/fullend/internal/contract"
	"github.com/geul-org/fullend/internal/projectconfig"
)

// generateIssueToken creates internal/auth/issue_token.go.
func generateIssueToken(authDir string, claims map[string]projectconfig.ClaimDef, fields []string, secretEnv, claimsHash string) error {
	var reqFields []string
	var mapClaims []string
	for _, field := range fields {
		def := claims[field]
		reqFields = append(reqFields, fmt.Sprintf("\t%s %s", field, def.GoType))
		mapClaims = append(mapClaims, fmt.Sprintf("\t\t\"%s\": req.%s,", def.Key, field))
	}

	envVar := "JWT_SECRET"
	if secretEnv != "" {
		envVar = secretEnv
	}

	src := fmt.Sprintf(`package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// @func issueToken
// @description claims 기반 JWT 액세스 토큰을 발급한다

type IssueTokenRequest struct {
%s
}

type IssueTokenResponse struct {
	AccessToken string
}

func IssueToken(req IssueTokenRequest) (IssueTokenResponse, error) {
	secret := os.Getenv(%q)
	if secret == "" {
		secret = "secret"
	}
	claims := jwt.MapClaims{
%s
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	return IssueTokenResponse{AccessToken: signed}, err
}
`, strings.Join(reqFields, "\n"), envVar, strings.Join(mapClaims, "\n"))

	d := &contract.Directive{Ownership: "gen", SSOT: "fullend.yaml", Contract: claimsHash}
	src = injectFileDirective(src, d)
	return os.WriteFile(filepath.Join(authDir, "issue_token.go"), []byte(src), 0644)
}
