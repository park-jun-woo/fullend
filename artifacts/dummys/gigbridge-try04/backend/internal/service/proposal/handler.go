package proposal

import (
	"database/sql"

	"github.com/example/gigbridge/internal/model"
)

// Handler handles requests for the proposal domain.
type Handler struct {
	DB *sql.DB
	GigModel model.GigModel
	ProposalModel model.ProposalModel
	TransactionModel model.TransactionModel
}
