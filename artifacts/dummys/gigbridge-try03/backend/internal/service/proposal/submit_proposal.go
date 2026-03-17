package proposal

import (
	"github.com/example/gigbridge/internal/model"
	"github.com/park-jun-woo/fullend/pkg/authz"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//fullend:gen ssot=service/proposal/submit_proposal.ssac contract=bcca69f
func (h *Handler) SubmitProposal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid path parameter"})
		return
	}

	currentUser := c.MustGet("currentUser").(*model.CurrentUser)

	var req struct {
		BidAmount int64 `json:"bid_amount" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	bidAmount := req.BidAmount

	tx, err := h.DB.BeginTx(c.Request.Context(), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "transaction failed"})
		return
	}
	defer tx.Rollback()

	if _, err = authz.Check(authz.CheckRequest{Action: "SubmitProposal", Resource: "proposal", UserID: currentUser.ID, Role: currentUser.Role}); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	proposal, err := h.ProposalModel.WithTx(tx).Create(id, currentUser.ID, bidAmount, "pending")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Proposal 생성 실패"})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "commit failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"proposal": proposal,
	})

}
