package template

import (
	"github.com/example/zenflow/internal/model"
	"github.com/park-jun-woo/fullend/pkg/authz"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//fullend:gen ssot=service/template/clone_template.ssac contract=7b8f4c9
func (h *Handler) CloneTemplate(c *gin.Context) {
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

	tmpl, err := h.TemplateModel.WithTx(tx).FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Template 조회 실패"})
		return
	}

	if tmpl == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	if _, err = authz.Check(authz.CheckRequest{Action: "CloneTemplate", Resource: "template", UserID: currentUser.ID, Role: currentUser.Role, ResourceID: tmpl.ID}); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	source, err := h.WorkflowModel.WithTx(tx).FindByID(tmpl.SourceWorkflowID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Workflow 조회 실패"})
		return
	}

	if source == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Source workflow not found"})
		return
	}

	newWf, err := h.WorkflowModel.WithTx(tx).Create(currentUser.OrgID, source.Title, source.TriggerEvent, "draft")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Workflow 생성 실패"})
		return
	}

	err = h.ActionModel.WithTx(tx).CopyToWorkflow(newWf.ID, source.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Action 수정 실패"})
		return
	}

	err = h.TemplateModel.WithTx(tx).IncrementCloneCount(tmpl.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Template 수정 실패"})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "commit failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"template_id": tmpl.ID,
		"workflow":    newWf,
	})

}
