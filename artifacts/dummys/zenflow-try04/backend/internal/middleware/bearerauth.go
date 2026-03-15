//fullend:gen ssot=fullend.yaml contract=78bda02
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/example/zenflow/internal/auth"
	"github.com/example/zenflow/internal/model"
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
			Email: out.Email,
			Role: out.Role,
			UserID: out.UserID,
		})
		c.Next()
	}
}
