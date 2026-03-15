package auth

import (
	"github.com/example/zenflow/internal/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

//fullend:gen ssot=service/auth/register.ssac contract=663b14a
func (h *Handler) Register(c *gin.Context) {
	var req struct {
		Password string `json:"password"`
		OrgName  string `json:"org_name"`
		Email    string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	password := req.Password
	orgName := req.OrgName
	email := req.Email

	tx, err := h.DB.BeginTx(c.Request.Context(), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "transaction failed"})
		return
	}
	defer tx.Rollback()

	hp, err := auth.HashPassword(auth.HashPasswordRequest{Password: password})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "호출 실패"})
		return
	}

	org, err := h.OrganizationModel.WithTx(tx).Create(orgName, "free")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Organization 생성 실패"})
		return
	}

	user, err := h.UserModel.WithTx(tx).Create(org.ID, email, hp.HashedPassword, "admin")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User 생성 실패"})
		return
	}

	token, err := auth.IssueToken(auth.IssueTokenRequest{Email: user.Email, ID: user.ID, OrgID: user.OrgID, Role: user.Role})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "호출 실패"})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "commit failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": token.AccessToken,
	})

}
