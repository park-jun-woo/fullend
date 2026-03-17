package gig

import (
	"github.com/example/gigbridge/internal/model"
	"github.com/park-jun-woo/fullend/pkg/authz"
	"github.com/example/gigbridge/internal/billing"
	"github.com/example/gigbridge/internal/mail"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/example/gigbridge/internal/states/gigstate"
	"strconv"
)

//fullend:gen ssot=service/gig/approve_work.ssac contract=37a3e8b
func (h *Handler) ApproveWork(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid path parameter"})
		return
	}

	currentUser := c.MustGet("currentUser").(*model.CurrentUser)

	tx, err := h.DB.BeginTx(c.Request.Context(), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "transaction failed"})
		return
	}
	defer tx.Rollback()

	gig, err := h.GigModel.WithTx(tx).FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gig 조회 실패"})
		return
	}

	if gig == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gig not found"})
		return
	}

	if _, err = authz.Check(authz.CheckRequest{Action: "ApproveWork", Resource: "gig", UserID: currentUser.ID, Role: currentUser.Role, ResourceID: gig.ID}); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	if err := gigstate.CanTransition(gigstate.Input{Status: gig.Status}, "ApproveWork"); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	_, err = billing.ReleaseFunds(billing.ReleaseFundsRequest{Amount: gig.Budget, FreelancerID: gig.FreelancerID, GigID: gig.ID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "호출 실패"})
		return
	}

	if _, err = mail.SendTemplateEmail(mail.SendTemplateEmailRequest{Subject: "Work Approved", TemplateName: "work_approved", To: currentUser.Email}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "호출 실패"})
		return
	}

	err = h.GigModel.WithTx(tx).UpdateStatus(gig.ID, "completed")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gig 수정 실패"})
		return
	}

	updated, err := h.GigModel.WithTx(tx).FindByID(gig.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gig 조회 실패"})
		return
	}

	if updated == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gig not found after update"})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "commit failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"gig": updated,
	})

}
