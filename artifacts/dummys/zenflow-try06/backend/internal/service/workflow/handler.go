package workflow

import (
	"database/sql"

	"github.com/example/zenflow/internal/model"
)

// Handler handles requests for the workflow domain.
type Handler struct {
	DB *sql.DB
	ActionModel model.ActionModel
	ExecutionLogModel model.ExecutionLogModel
	OrganizationModel model.OrganizationModel
	WorkflowModel model.WorkflowModel
}
