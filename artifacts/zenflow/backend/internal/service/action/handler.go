package action

import (
	"database/sql"

	"github.com/geul-org/zenflow/internal/model"
)

// Handler handles requests for the action domain.
type Handler struct {
	DB *sql.DB
	ActionModel model.ActionModel
	WorkflowModel model.WorkflowModel
}
