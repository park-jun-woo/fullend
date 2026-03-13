package workflow

import (
	"github.com/zenflow/zenflow/internal/model"
	"github.com/geul-org/fullend/pkg/authz"
	"github.com/gin-gonic/gin"
	"net/http"
)

//fullend:gen ssot=service/workflow/list_workflows.ssac contract=f65b04a
func (h *Handler) ListWorkflows(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(*model.CurrentUser)

	if _, err := authz.Check(authz.CheckRequest{Action: "ListWorkflows", Resource: "workflow", UserID: currentUser.ID, Role: currentUser.Role}); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	user, err := h.UserModel.FindByID(currentUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User 조회 실패"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	workflows, err := h.WorkflowModel.List(user.OrgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Workflow 조회 실패"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"workflows": workflows,
	})

}
