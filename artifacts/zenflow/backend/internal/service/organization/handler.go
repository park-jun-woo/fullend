package organization

import (
	"database/sql"

	"github.com/geul-org/zenflow/internal/model"
)

// Handler handles requests for the organization domain.
type Handler struct {
	DB *sql.DB
	OrganizationModel model.OrganizationModel
}
