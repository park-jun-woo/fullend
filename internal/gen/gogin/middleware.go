package gogin

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/geul-org/fullend/internal/contract"
	"github.com/geul-org/fullend/internal/projectconfig"
)

// generateAuthPackage creates internal/auth/ with claims-based JWT functions and reexport.
func generateAuthPackage(intDir, modulePath string, claims map[string]projectconfig.ClaimDef, secretEnv string) error {
	authDir := filepath.Join(intDir, "auth")
	if err := os.MkdirAll(authDir, 0755); err != nil {
		return err
	}

	fields := sortedClaimFields(claims)
	claimsHash := HashClaimDefs(claims)

	if err := generateIssueToken(authDir, claims, fields, secretEnv, claimsHash); err != nil {
		return err
	}
	if err := generateVerifyToken(authDir, claims, fields, claimsHash); err != nil {
		return err
	}
	if err := generateRefreshToken(authDir, claims, fields, claimsHash); err != nil {
		return err
	}
	if err := generateReexport(authDir, claimsHash); err != nil {
		return err
	}

	return nil
}

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

// generateRefreshToken creates internal/auth/refresh_token.go.
func generateRefreshToken(authDir string, claims map[string]projectconfig.ClaimDef, fields []string, claimsHash string) error {
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

// generateReexport creates internal/auth/reexport.go that re-exports pkg/auth utilities.
func generateReexport(authDir string, claimsHash string) error {
	src := `package auth

import pkgauth "github.com/geul-org/fullend/pkg/auth"

// Re-export pkg/auth utilities for unified import.
var HashPassword = pkgauth.HashPassword
var VerifyPassword = pkgauth.VerifyPassword
var GenerateResetToken = pkgauth.GenerateResetToken

type HashPasswordRequest = pkgauth.HashPasswordRequest
type HashPasswordResponse = pkgauth.HashPasswordResponse
type VerifyPasswordRequest = pkgauth.VerifyPasswordRequest
type VerifyPasswordResponse = pkgauth.VerifyPasswordResponse
type GenerateResetTokenRequest = pkgauth.GenerateResetTokenRequest
type GenerateResetTokenResponse = pkgauth.GenerateResetTokenResponse
`
	d := &contract.Directive{Ownership: "gen", SSOT: "fullend.yaml", Contract: claimsHash}
	src = injectFileDirective(src, d)
	return os.WriteFile(filepath.Join(authDir, "reexport.go"), []byte(src), 0644)
}

// generateMiddleware creates internal/middleware/bearerauth.go with claims-based CurrentUser mapping.
func generateMiddleware(intDir, modulePath string, claims map[string]projectconfig.ClaimDef) error {
	mwDir := filepath.Join(intDir, "middleware")
	if err := os.MkdirAll(mwDir, 0755); err != nil {
		return err
	}

	fields := sortedClaimFields(claims)
	claimsHash := HashClaimDefs(claims)

	var assignments []string
	for _, field := range fields {
		assignments = append(assignments, fmt.Sprintf("\t\t\t%s: out.%s,", field, field))
	}
	assignBlock := strings.Join(assignments, "\n")

	src := fmt.Sprintf(`package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"%s/internal/auth"
	"%s/internal/model"
)

// BearerAuth returns a gin middleware that validates the Authorization header.
func BearerAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		token := strings.TrimPrefix(header, "Bearer ")
		out, err := auth.VerifyToken(auth.VerifyTokenRequest{Token: token, Secret: secret})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Set("currentUser", &model.CurrentUser{
%s
		})
		c.Next()
	}
}
`, modulePath, modulePath, assignBlock)

	d := &contract.Directive{Ownership: "gen", SSOT: "fullend.yaml", Contract: claimsHash}
	src = injectFileDirective(src, d)
	return os.WriteFile(filepath.Join(mwDir, "bearerauth.go"), []byte(src), 0644)
}

// sortedClaimFields returns claim field names in sorted order.
func sortedClaimFields(claims map[string]projectconfig.ClaimDef) []string {
	var fields []string
	for field := range claims {
		fields = append(fields, field)
	}
	sort.Strings(fields)
	return fields
}

// HashClaimDefs computes a contract hash for ClaimDef claims.
func HashClaimDefs(claims map[string]projectconfig.ClaimDef) string {
	fields := sortedClaimFields(claims)
	var parts []string
	for _, f := range fields {
		def := claims[f]
		parts = append(parts, f+":"+def.Key+":"+def.GoType)
	}
	return contract.Hash7(strings.Join(parts, ","))
}

// claimExtractLine generates the JWT MapClaims extraction line for VerifyToken.
func claimExtractLine(field string, def projectconfig.ClaimDef) string {
	varName := lcFirst(field) + "Raw"
	switch def.GoType {
	case "int64":
		return fmt.Sprintf("\t%s, _ := claims[\"%s\"].(float64)", varName, def.Key)
	case "bool":
		return fmt.Sprintf("\t%s, _ := claims[\"%s\"].(bool)", varName, def.Key)
	default: // string
		return fmt.Sprintf("\t%s, _ := claims[\"%s\"].(string)", varName, def.Key)
	}
}

// resultAssignLine generates the struct field assignment for VerifyToken result.
func resultAssignLine(field string, def projectconfig.ClaimDef) string {
	varName := lcFirst(field) + "Raw"
	switch def.GoType {
	case "int64":
		return fmt.Sprintf("\t\t%s: int64(%s),", field, varName)
	default:
		return fmt.Sprintf("\t\t%s: %s,", field, varName)
	}
}
