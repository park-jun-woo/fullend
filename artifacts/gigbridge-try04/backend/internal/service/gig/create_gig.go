package gig

import (
	"github.com/example/gigbridge/internal/model"
	"github.com/park-jun-woo/fullend/pkg/authz"
	"github.com/gin-gonic/gin"
	"net/http"
)

//fullend:gen ssot=service/gig/create_gig.ssac contract=10dc40c
func (h *Handler) CreateGig(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(*model.CurrentUser)

	var req struct {
		Title       string `json:"title" binding:"required,max=255"`
		Description string `json:"description" binding:"required"`
		Budget      int64  `json:"budget" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	title := req.Title
	description := req.Description
	budget := req.Budget

	tx, err := h.DB.BeginTx(c.Request.Context(), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "transaction failed"})
		return
	}
	defer tx.Rollback()

	if _, err = authz.Check(authz.CheckRequest{Action: "CreateGig", Resource: "gig", UserID: currentUser.ID, Role: currentUser.Role, ResourceID: currentUser.ID}); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	gig, err := h.GigModel.WithTx(tx).Create(currentUser.ID, title, description, budget, "draft")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gig 생성 실패"})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "commit failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"gig": gig,
	})

}
