package webhook

import (
	"context"
	"fmt"
	"github.com/example/zenflow/internal/webhook"
)

type WorkflowExecutedMessage struct {
	WorkflowID int64
	OrgID      int64
	Status     string
}

//fullend:gen ssot=service/webhook/on_workflow_executed.ssac contract=e69d223
func (h *Handler) OnWorkflowExecuted(ctx context.Context, message WorkflowExecutedMessage) error {
	_, err := h.WebhookModel.ListByOrgIDAndEventType(message.OrgID, "workflow.executed")
	if err != nil {
		return fmt.Errorf("Webhook 조회 실패: %w", err)
	}

	if _, err = webhook.Deliver(webhook.DeliverRequest{OrgID: message.OrgID, Status: message.Status, WorkflowID: message.WorkflowID}); err != nil {
		return fmt.Errorf("호출 실패: %w", err)
	}

	return nil
}
