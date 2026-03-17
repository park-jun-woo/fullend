package proposal

import (
	"github.com/example/gigbridge/internal/model"
	"github.com/park-jun-woo/fullend/pkg/authz"
	"github.com/example/gigbridge/internal/billing"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/example/gigbridge/internal/states/gigstate"
	"github.com/example/gigbridge/internal/states/proposalstate"
	"strconv"
)

//fullend:gen ssot=service/proposal/accept_proposal.ssac contract=aa2292b
func (h *Handler) AcceptProposal(c *gin.Context) {
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

	proposal, err := h.ProposalModel.WithTx(tx).FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Proposal 조회 실패"})
		return
	}

	if proposal == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Proposal not found"})
		return
	}

	gig, err := h.GigModel.WithTx(tx).FindByID(proposal.GigID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gig 조회 실패"})
		return
	}

	if gig == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gig not found"})
		return
	}

	if _, err = authz.Check(authz.CheckRequest{Action: "AcceptProposal", Resource: "gig", UserID: currentUser.ID, Role: currentUser.Role, ResourceID: gig.ID}); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	if err := proposalstate.CanTransition(proposalstate.Input{Status: proposal.Status}, "AcceptProposal"); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	if err := gigstate.CanTransition(gigstate.Input{Status: gig.Status}, "AcceptProposal"); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	_, err = billing.HoldEscrow(billing.HoldEscrowRequest{Amount: gig.Budget, ClientID: gig.ClientID, GigID: gig.ID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "호출 실패"})
		return
	}

	err = h.ProposalModel.WithTx(tx).UpdateStatus(proposal.ID, "accepted")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Proposal 수정 실패"})
		return
	}

	err = h.GigModel.WithTx(tx).UpdateStatus(gig.ID, "in_progress")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gig 수정 실패"})
		return
	}

	err = h.GigModel.WithTx(tx).UpdateFreelancer(gig.ID, proposal.FreelancerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gig 수정 실패"})
		return
	}

	updatedProposal, err := h.ProposalModel.WithTx(tx).FindByID(proposal.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Proposal 조회 실패"})
		return
	}

	if updatedProposal == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Proposal not found after update"})
		return
	}

	updatedGig, err := h.GigModel.WithTx(tx).FindByID(gig.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gig 조회 실패"})
		return
	}

	if updatedGig == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gig not found after update"})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "commit failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"gig":      updatedGig,
		"proposal": updatedProposal,
	})

}
