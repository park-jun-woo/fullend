package webhook

import (
	"github.com/example/zenflow/internal/model"
	"github.com/park-jun-woo/fullend/pkg/authz"
	"github.com/gin-gonic/gin"
	"net/http"
)

//fullend:gen ssot=service/webhook/list_webhooks.ssac contract=ee381b2
func (h *Handler) ListWebhooks(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(*model.CurrentUser)

	if _, err := authz.Check(authz.CheckRequest{Action: "ListWebhooks", Resource: "webhook", UserID: currentUser.ID, Role: currentUser.Role, ResourceID: currentUser.ID}); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	webhooks, err := h.WebhookModel.ListByOrgID(currentUser.OrgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Webhook 조회 실패"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"webhooks": webhooks,
	})

}
