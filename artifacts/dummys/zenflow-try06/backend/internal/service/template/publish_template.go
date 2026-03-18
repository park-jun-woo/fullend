package template

import (
	"github.com/example/zenflow/internal/model"
	"github.com/park-jun-woo/fullend/pkg/authz"
	"github.com/gin-gonic/gin"
	"net/http"
)

//fullend:gen ssot=service/template/publish_template.ssac contract=1159557
func (h *Handler) PublishTemplate(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(*model.CurrentUser)

	var req struct {
		SourceWorkflowID int64  `json:"source_workflow_id" binding:"required"`
		Category         string `json:"category" binding:"required,max=100"`
		Title            string `json:"title" binding:"required,max=255"`
		Description      string `json:"description" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sourceWorkflowID := req.SourceWorkflowID
	category := req.Category
	title := req.Title
	description := req.Description

	tx, err := h.DB.BeginTx(c.Request.Context(), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "transaction failed"})
		return
	}
	defer tx.Rollback()

	wf, err := h.WorkflowModel.WithTx(tx).FindByIDAndOrgID(sourceWorkflowID, currentUser.OrgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Workflow 조회 실패"})
		return
	}

	if wf == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
		return
	}

	if _, err = authz.Check(authz.CheckRequest{Action: "PublishTemplate", Resource: "template", UserID: currentUser.ID, Role: currentUser.Role, ResourceID: currentUser.ID}); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	existing, err := h.TemplateModel.WithTx(tx).FindBySourceWorkflowID(wf.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Template 조회 실패"})
		return
	}

	if existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Template already published for this workflow"})
		return
	}

	tmpl, err := h.TemplateModel.WithTx(tx).Create(wf.ID, currentUser.OrgID, title, description, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Template 생성 실패"})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "commit failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"template": tmpl,
	})

}
