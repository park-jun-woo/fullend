package workflow

import (
	"database/sql"

	"github.com/geul-org/zenflow/internal/model"
)

// Handler handles requests for the workflow domain.
type Handler struct {
	DB *sql.DB
	ExecutionModel model.ExecutionModel
	UserModel model.UserModel
	WorkflowModel model.WorkflowModel
}
