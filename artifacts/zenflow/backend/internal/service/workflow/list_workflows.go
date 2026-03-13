package workflow

import (
	"github.com/geul-org/zenflow/internal/model"
	"github.com/geul-org/fullend/pkg/authz"
	"github.com/gin-gonic/gin"
	"net/http"
)

//fullend:gen ssot=service/workflow/list_workflows.ssac contract=090370f
func (h *Handler) ListWorkflows(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(*model.CurrentUser)

	me, err := h.UserModel.FindByID(currentUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User 조회 실패"})
		return
	}

	if me == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if _, err = authz.Check(authz.CheckRequest{Action: "ListWorkflows", Resource: "workflow", UserID: currentUser.ID, Role: currentUser.Role, ResourceID: me.OrgID}); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	workflows, err := h.WorkflowModel.ListByOrgID(me.OrgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Workflow 조회 실패"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"workflows": workflows,
	})

}
