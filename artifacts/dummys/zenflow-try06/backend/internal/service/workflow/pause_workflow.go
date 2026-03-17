package workflow

import (
	"github.com/example/zenflow/internal/model"
	"github.com/park-jun-woo/fullend/pkg/authz"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/example/zenflow/internal/states/workflowstate"
	"strconv"
)

//fullend:gen ssot=service/workflow/pause_workflow.ssac contract=a914131
func (h *Handler) PauseWorkflow(c *gin.Context) {
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

	wf, err := h.WorkflowModel.WithTx(tx).FindByIDAndOrgID(id, currentUser.OrgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Workflow 조회 실패"})
		return
	}

	if wf == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
		return
	}

	if _, err = authz.Check(authz.CheckRequest{Action: "PauseWorkflow", Resource: "workflow", UserID: currentUser.ID, Role: currentUser.Role, ResourceID: wf.ID}); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	if err := workflowstate.CanTransition(workflowstate.Input{Status: wf.Status}, "PauseWorkflow"); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	err = h.WorkflowModel.WithTx(tx).UpdateStatus(wf.ID, "paused")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Workflow 수정 실패"})
		return
	}

	updated, err := h.WorkflowModel.WithTx(tx).FindByID(wf.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Workflow 조회 실패"})
		return
	}

	if updated == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "commit failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"workflow": updated,
	})

}
