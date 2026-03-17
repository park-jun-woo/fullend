package webhook

import (
	"github.com/example/zenflow/internal/model"
	"github.com/park-jun-woo/fullend/pkg/authz"
	"github.com/gin-gonic/gin"
	"net/http"
)

//fullend:gen ssot=service/webhook/create_webhook.ssac contract=b7ff8be
func (h *Handler) CreateWebhook(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(*model.CurrentUser)

	var req struct {
		EventType string `json:"event_type" binding:"required,max=100"`
		URL       string `json:"url" binding:"required,max=2048"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	eventType := req.EventType
	url := req.URL

	tx, err := h.DB.BeginTx(c.Request.Context(), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "transaction failed"})
		return
	}
	defer tx.Rollback()

	if _, err = authz.Check(authz.CheckRequest{Action: "CreateWebhook", Resource: "webhook", UserID: currentUser.ID, Role: currentUser.Role, ResourceID: currentUser.ID}); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	wh, err := h.WebhookModel.WithTx(tx).Create(currentUser.OrgID, url, eventType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Webhook 생성 실패"})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "commit failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"webhook": wh,
	})

}
