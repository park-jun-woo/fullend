package log

import (
	"github.com/example/zenflow/internal/model"
	"github.com/park-jun-woo/fullend/pkg/authz"
	"github.com/gin-gonic/gin"
	"net/http"
)

//fullend:gen ssot=service/log/list_execution_logs.ssac contract=aca9bbb
func (h *Handler) ListExecutionLogs(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(*model.CurrentUser)

	if _, err := authz.Check(authz.CheckRequest{Action: "ListExecutionLogs", Resource: "execution_log", UserID: currentUser.ID, Role: currentUser.Role}); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	logs, err := h.ExecutionLogModel.ListByOrg(currentUser.OrgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ExecutionLog 조회 실패"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"execution_logs": logs,
	})

}
