package workflow

import (
	"github.com/example/zenflow/internal/model"
	"github.com/park-jun-woo/fullend/pkg/authz"
	"github.com/example/zenflow/internal/billing"
	"github.com/example/zenflow/internal/worker"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/example/zenflow/internal/states/workflowstate"
	"strconv"
)

//fullend:gen ssot=service/workflow/execute_workflow.ssac contract=8eac9cb
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

	if _, err = authz.Check(authz.CheckRequest{Action: "ExecuteWorkflow", Resource: "workflow", UserID: currentUser.UserID, Role: currentUser.Role, ResourceID: id}); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	wf, err := h.WorkflowModel.WithTx(tx).FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Workflow 조회 실패"})
		return
	}

	if wf == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
		return
	}

	if err := workflowstate.CanTransition(workflowstate.Input{Status: wf.Status}, "ExecuteWorkflow"); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	org, err := h.OrganizationModel.WithTx(tx).FindByID(wf.OrgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Organization 조회 실패"})
		return
	}

	if org == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	if _, err = billing.ValidateCredits(billing.ValidateCreditsRequest{CreditsBalance: org.CreditsBalance}); err != nil {
		c.JSON(http.StatusPaymentRequired, gin.H{"error": "호출 실패"})
		return
	}

	_, err = h.ActionModel.WithTx(tx).ListByWorkflowID(wf.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Action 조회 실패"})
		return
	}

	_, err = worker.ProcessAction(worker.ProcessActionRequest{ActionType: wf.TriggerEvent, PayloadTemplate: wf.Title})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "호출 실패"})
		return
	}

	err = h.OrganizationModel.WithTx(tx).DeductCredit(org.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Organization 수정 실패"})
		return
	}

	log, err := h.ExecutionLogModel.WithTx(tx).Create(wf.ID, wf.OrgID, "completed")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ExecutionLog 생성 실패"})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "commit failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"execution_log": log,
	})

}
