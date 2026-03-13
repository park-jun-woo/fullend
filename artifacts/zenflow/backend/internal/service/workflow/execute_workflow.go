package workflow

import (
	"github.com/zenflow/zenflow/internal/model"
	"github.com/geul-org/fullend/pkg/authz"
	"github.com/gin-gonic/gin"
	"github.com/zenflow/zenflow/internal/worker"
	"net/http"
	"github.com/zenflow/zenflow/internal/states/workflowstate"
	"strconv"
)

//fullend:gen ssot=service/workflow/execute_workflow.ssac contract=ce776a0
func (h *Handler) ExecuteWorkflow(c *gin.Context) {
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

	if _, err = authz.Check(authz.CheckRequest{Action: "ExecuteWorkflow", Resource: "workflow", UserID: currentUser.ID, Role: currentUser.Role}); err != nil {
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

	if err := workflowstate.CanTransition(workflowstate.Input{Status: workflow.Status}, "ExecuteWorkflow"); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	_, err = h.ActionModel.WithTx(tx).ListByWorkflow(workflow.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Action 조회 실패"})
		return
	}

	result, err := worker.ProcessAction(worker.ProcessActionRequest{WorkflowID: workflow.ID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "호출 실패"})
		return
	}

	org, err := h.OrganizationModel.WithTx(tx).FindByID(user.OrgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Organization 조회 실패"})
		return
	}

	if org == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	err = h.OrganizationModel.WithTx(tx).UpdateCredits(org.ID, result.ProcessedCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Organization 수정 실패"})
		return
	}

	execution, err := h.ExecutionModel.WithTx(tx).Create(workflow.ID, org.ID, "completed", result.ProcessedCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Execution 생성 실패"})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "commit failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"execution": execution,
	})

}
