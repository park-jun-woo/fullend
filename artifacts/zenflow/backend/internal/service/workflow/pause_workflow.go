package workflow

import (
	"github.com/zenflow/zenflow/internal/model"
	"github.com/geul-org/fullend/pkg/authz"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/zenflow/zenflow/internal/states/workflowstate"
	"strconv"
)

//fullend:gen ssot=service/workflow/pause_workflow.ssac contract=d6feae2
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

	if _, err = authz.Check(authz.CheckRequest{Action: "PauseWorkflow", Resource: "workflow", UserID: currentUser.ID, Role: currentUser.Role}); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	user, err := h.UserModel.WithTx(tx).FindByID(currentUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User 조회 실패"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	workflow, err := h.WorkflowModel.WithTx(tx).FindByIDAndOrg(id, user.OrgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Workflow 조회 실패"})
		return
	}

	if workflow == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
		return
	}

	if err := workflowstate.CanTransition(workflowstate.Input{Status: workflow.Status}, "PauseWorkflow"); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	err = h.WorkflowModel.WithTx(tx).UpdateStatus(workflow.ID, "paused")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Workflow 수정 실패"})
		return
	}

	updated, err := h.WorkflowModel.WithTx(tx).FindByID(workflow.ID)
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
