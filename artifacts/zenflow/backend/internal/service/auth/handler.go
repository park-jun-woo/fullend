package auth

import (
	"database/sql"

	"github.com/geul-org/zenflow/internal/model"
)

// Handler handles requests for the auth domain.
type Handler struct {
	DB *sql.DB
	UserModel model.UserModel
	JWTSecret string
}
