package workflow

import (
	"github.com/example/zenflow/internal/model"
	"github.com/park-jun-woo/fullend/pkg/authz"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//fullend:gen ssot=service/workflow/list_workflow_versions.ssac contract=068b316
func (h *Handler) ListWorkflowVersions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid path parameter"})
		return
	}

	currentUser := c.MustGet("currentUser").(*model.CurrentUser)

	wf, err := h.WorkflowModel.FindByIDAndOrgID(id, currentUser.OrgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Workflow 조회 실패"})
		return
	}

	if wf == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
		return
	}

	if _, err = authz.Check(authz.CheckRequest{Action: "ListWorkflowVersions", Resource: "workflow", UserID: currentUser.ID, Role: currentUser.Role, ResourceID: wf.ID}); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	versions, err := h.WorkflowModel.ListVersions(wf.ID, currentUser.OrgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Workflow 조회 실패"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"versions": versions,
	})

}
