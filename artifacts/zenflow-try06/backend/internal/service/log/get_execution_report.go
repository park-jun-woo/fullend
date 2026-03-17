package log

import (
	"github.com/example/zenflow/internal/model"
	"github.com/park-jun-woo/fullend/pkg/authz"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//fullend:gen ssot=service/log/get_execution_report.ssac contract=7adf2ec
func (h *Handler) GetExecutionReport(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid path parameter"})
		return
	}

	currentUser := c.MustGet("currentUser").(*model.CurrentUser)

	execLog, err := h.ExecutionLogModel.FindByIDAndOrgID(id, currentUser.OrgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ExecutionLog 조회 실패"})
		return
	}

	if execLog == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Execution log not found"})
		return
	}

	if _, err = authz.Check(authz.CheckRequest{Action: "GetExecutionReport", Resource: "workflow", UserID: currentUser.ID, Role: currentUser.Role, ResourceID: execLog.ID}); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"credits_spent": execLog.CreditsSpent,
		"report_key":    execLog.ReportKey,
		"status":        execLog.Status,
	})

}
