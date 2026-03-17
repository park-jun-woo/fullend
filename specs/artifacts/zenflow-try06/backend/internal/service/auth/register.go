package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/example/zenflow/internal/auth"
	"net/http"
)

//fullend:gen ssot=service/auth/register.ssac contract=3ded654
func (h *Handler) Register(c *gin.Context) {
	var req struct {
		OrgID    int64  `json:"org_id" binding:"required"`
		Password string `json:"password" binding:"required,min=8"`
		Email    string `json:"email" binding:"required,email,max=255"`
		Role     string `json:"role" binding:"required,max=50"`
		Name     string `json:"name" binding:"required,max=255"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	orgID := req.OrgID
	password := req.Password
	email := req.Email
	role := req.Role
	name := req.Name

	tx, err := h.DB.BeginTx(c.Request.Context(), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "transaction failed"})
		return
	}
	defer tx.Rollback()

	org, err := h.OrganizationModel.WithTx(tx).FindByID(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Organization 조회 실패"})
		return
	}

	if org == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	hp, err := auth.HashPassword(auth.HashPasswordRequest{Password: password})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "호출 실패"})
		return
	}

	user, err := h.UserModel.WithTx(tx).Create(orgID, email, hp.HashedPassword, role, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User 생성 실패"})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "commit failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": user,
	})

}
