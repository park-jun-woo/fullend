package workflow

import (
	"github.com/example/zenflow/internal/model"
	"github.com/park-jun-woo/fullend/pkg/authz"
	"github.com/example/zenflow/internal/billing"
	"github.com/example/zenflow/internal/report"
	"github.com/example/zenflow/internal/webhook"
	"github.com/example/zenflow/internal/worker"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/park-jun-woo/fullend/pkg/queue"
	"github.com/example/zenflow/internal/states/workflowstate"
	"strconv"
)

//fullend:gen ssot=service/workflow/execute_with_report.ssac contract=b119aec
func (h *Handler) ExecuteWithReport(c *gin.Context) {
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

	if _, err = authz.Check(authz.CheckRequest{Action: "ExecuteWithReport", Resource: "workflow", UserID: currentUser.ID, Role: currentUser.Role, ResourceID: wf.ID}); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	if err := workflowstate.CanTransition(workflowstate.Input{Status: wf.Status}, "ExecuteWorkflow"); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	_, err = h.ActionModel.WithTx(tx).ListByWorkflowID(wf.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Action 조회 실패"})
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

	_, err = billing.CheckCredits(billing.CheckCreditsRequest{Balance: org.CreditsBalance, OrgID: org.ID})
	if err != nil {
		c.JSON(http.StatusPaymentRequired, gin.H{"error": "호출 실패"})
		return
	}

	result, err := worker.ProcessActions(worker.ProcessActionsRequest{ActionCount: org.CreditsBalance, WorkflowID: wf.ID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "호출 실패"})
		return
	}

	deduction, err := billing.DeductCredit(billing.DeductCreditRequest{Amount: result.CreditsUsed, CurrentBalance: org.CreditsBalance, OrgID: wf.OrgID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "호출 실패"})
		return
	}

	err = h.OrganizationModel.WithTx(tx).UpdateCredits(wf.OrgID, deduction.NewBalance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Organization 수정 실패"})
		return
	}

	rpt, err := report.GenerateReport(report.GenerateReportRequest{Status: "completed", WorkflowID: wf.ID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "호출 실패"})
		return
	}

	execLog, err := h.ExecutionLogModel.WithTx(tx).CreateWithReport(wf.ID, wf.OrgID, "completed", result.CreditsUsed, rpt.ReportKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ExecutionLog 생성 실패"})
		return
	}

	err = queue.Publish(c.Request.Context(), "workflow.executed", map[string]any{
		"OrgID":      wf.OrgID,
		"Status":     "completed",
		"WorkflowID": wf.ID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "처리 실패"})
		return
	}

	_, err = webhook.Deliver(webhook.DeliverRequest{OrgID: wf.OrgID, Status: "completed", WorkflowID: wf.ID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "호출 실패"})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "commit failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"credits_remaining": deduction.NewBalance,
		"log":               execLog,
		"report_key":        rpt.ReportKey,
	})

}
