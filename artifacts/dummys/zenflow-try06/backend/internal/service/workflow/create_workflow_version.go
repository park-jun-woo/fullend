package workflow

import (
	"github.com/example/zenflow/internal/model"
	"github.com/park-jun-woo/fullend/pkg/authz"
	"github.com/example/zenflow/internal/workflow"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//fullend:gen ssot=service/workflow/create_workflow_version.ssac contract=6c7a321
func (h *Handler) CreateWorkflowVersion(c *gin.Context) {
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

	source, err := h.WorkflowModel.WithTx(tx).FindByIDAndOrgID(id, currentUser.OrgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Workflow 조회 실패"})
		return
	}

	if source == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
		return
	}

	if _, err = authz.Check(authz.CheckRequest{Action: "CreateWorkflowVersion", Resource: "workflow", UserID: currentUser.ID, Role: currentUser.Role, ResourceID: source.ID}); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	root, err := workflow.ResolveRootID(workflow.ResolveRootIDRequest{RootWorkflowID: source.RootWorkflowID, WorkflowID: source.ID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "호출 실패"})
		return
	}

	ver, err := workflow.NextVersion(workflow.NextVersionRequest{CurrentVersion: source.Version})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "호출 실패"})
		return
	}

	newWf, err := h.WorkflowModel.WithTx(tx).CreateVersion(source.OrgID, source.Title, source.TriggerEvent, "draft", ver.NextVersion, root.ResolvedRootID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Workflow 생성 실패"})
		return
	}

	err = h.ActionModel.WithTx(tx).CopyToWorkflow(newWf.ID, source.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Action 수정 실패"})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "commit failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"workflow": newWf,
	})

}
