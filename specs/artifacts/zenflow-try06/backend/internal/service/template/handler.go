package template

import (
	"database/sql"

	"github.com/example/zenflow/internal/model"
)

// Handler handles requests for the template domain.
type Handler struct {
	DB *sql.DB
	ActionModel model.ActionModel
	OrganizationModel model.OrganizationModel
	TemplateModel model.TemplateModel
	WorkflowModel model.WorkflowModel
}
