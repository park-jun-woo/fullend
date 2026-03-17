package workflow

import (
	"github.com/example/zenflow/internal/model"
	"github.com/park-jun-woo/fullend/pkg/authz"
	"github.com/gin-gonic/gin"
	"net/http"
)

//fullend:gen ssot=service/workflow/list_workflows.ssac contract=a4a7a61
func (h *Handler) ListWorkflows(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(*model.CurrentUser)

	if _, err := authz.Check(authz.CheckRequest{Action: "ListWorkflows", Resource: "workflow", UserID: currentUser.UserID, Role: currentUser.Role}); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	viewer, err := h.UserModel.FindByID(currentUser.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User 조회 실패"})
		return
	}

	workflows, err := h.WorkflowModel.ListByOrgID(viewer.OrgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Workflow 조회 실패"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"workflows": workflows,
	})

}
