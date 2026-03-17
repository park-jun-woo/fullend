package webhook

import (
	"github.com/example/zenflow/internal/model"
	"github.com/park-jun-woo/fullend/pkg/authz"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//fullend:gen ssot=service/webhook/delete_webhook.ssac contract=0727ca0
func (h *Handler) DeleteWebhook(c *gin.Context) {
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

	wh, err := h.WebhookModel.WithTx(tx).FindByIDAndOrgID(id, currentUser.OrgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Webhook 조회 실패"})
		return
	}

	if wh == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Webhook not found"})
		return
	}

	if _, err = authz.Check(authz.CheckRequest{Action: "DeleteWebhook", Resource: "webhook", UserID: currentUser.ID, Role: currentUser.Role, ResourceID: wh.ID}); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	err = h.WebhookModel.WithTx(tx).Delete(wh.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Webhook 삭제 실패"})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "commit failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": wh.ID,
	})

}
