package schedule

import (
	"github.com/example/zenflow/internal/model"
	"github.com/park-jun-woo/fullend/pkg/authz"
	"github.com/example/zenflow/internal/schedule"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//fullend:gen ssot=service/schedule/set_schedule.ssac contract=8ae4b4c
func (h *Handler) SetSchedule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid path parameter"})
		return
	}

	currentUser := c.MustGet("currentUser").(*model.CurrentUser)

	var req struct {
		Cron string `json:"cron" binding:"required,max=100"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cron := req.Cron

	wf, err := h.WorkflowModel.FindByIDAndOrgID(id, currentUser.OrgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Workflow 조회 실패"})
		return
	}

	if wf == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
		return
	}

	if _, err = authz.Check(authz.CheckRequest{Action: "SetSchedule", Resource: "workflow", UserID: currentUser.ID, Role: currentUser.Role, ResourceID: wf.ID}); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	result, err := schedule.SetSchedule(schedule.SetScheduleRequest{Cron: cron, WorkflowID: wf.ID})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "호출 실패"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cron":     result.Cron,
		"next_run": result.NextRun,
	})

}
