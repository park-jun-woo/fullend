//ff:func feature=gen-gogin type=generator control=iteration
//ff:what creates internal/middleware/bearerauth.go with claims-based CurrentUser mapping

package gogin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/geul-org/fullend/internal/contract"
	"github.com/geul-org/fullend/internal/projectconfig"
)

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
