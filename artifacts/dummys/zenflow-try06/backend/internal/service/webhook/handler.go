package webhook

import (
	"database/sql"

	"github.com/example/zenflow/internal/model"
)

// Handler handles requests for the webhook domain.
type Handler struct {
	DB *sql.DB
	WebhookModel model.WebhookModel
}
