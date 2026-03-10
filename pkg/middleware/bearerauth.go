package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/geul-org/fullend/pkg/auth"
)

// CurrentUser represents the authenticated user extracted from JWT.
type CurrentUser struct {
	ID    int64
	Email string
	Role  string
}

// BearerAuth returns a gin middleware that extracts the current user from the Authorization header.
// It does NOT abort — authorize sequences in handlers check permissions.
// Name matches OpenAPI securitySchemes key "bearerAuth".
func BearerAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.Set("currentUser", &CurrentUser{})
			c.Next()
			return
		}
		token := strings.TrimPrefix(header, "Bearer ")
		out, err := auth.VerifyToken(auth.VerifyTokenRequest{Token: token, Secret: secret})
		if err != nil {
			c.Set("currentUser", &CurrentUser{})
			c.Next()
			return
		}
		c.Set("currentUser", &CurrentUser{
			ID:    out.UserID,
			Email: out.Email,
			Role:  out.Role,
		})
		c.Next()
	}
}
