package auth

import (
	"database/sql"

	"github.com/example/zenflow/internal/model"
)

// Handler handles requests for the auth domain.
type Handler struct {
	DB *sql.DB
	OrganizationModel model.OrganizationModel
	UserModel model.UserModel
	JWTSecret string
}
