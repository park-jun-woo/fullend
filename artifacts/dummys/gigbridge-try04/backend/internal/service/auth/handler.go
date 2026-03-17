package auth

import (
	"database/sql"

	"github.com/example/gigbridge/internal/model"
)

// Handler handles requests for the auth domain.
type Handler struct {
	DB *sql.DB
	UserModel model.UserModel
	JWTSecret string
}
