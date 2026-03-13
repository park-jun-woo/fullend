package workflow

import (
	"github.com/geul-org/zenflow/internal/model"
	"github.com/geul-org/fullend/pkg/authz"
	"github.com/gin-gonic/gin"
	"net/http"
)

//fullend:gen ssot=service/workflow/create_workflow.ssac contract=51b8c5b
func (h *Handler) CreateWorkflow(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(*model.CurrentUser)

	var req struct {
		Title        string `json:"title"`
		TriggerEvent string `json:"trigger_event"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	title := req.Title
	triggerEvent := req.TriggerEvent

	tx, err := h.DB.BeginTx(c.Request.Context(), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "transaction failed"})
		return
	}
	defer tx.Rollback()

	me, err := h.UserModel.WithTx(tx).FindByID(currentUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User 조회 실패"})
		return
	}

	if me == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if _, err = authz.Check(authz.CheckRequest{Action: "CreateWorkflow", Resource: "workflow", UserID: currentUser.ID, Role: currentUser.Role, ResourceID: me.OrgID}); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	wf, err := h.WorkflowModel.WithTx(tx).Create(me.OrgID, title, triggerEvent, "draft")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Workflow 생성 실패"})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "commit failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"workflow": wf,
	})

}
