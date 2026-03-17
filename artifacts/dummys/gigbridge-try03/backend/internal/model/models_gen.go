package model

import (
	"database/sql"

	"github.com/park-jun-woo/fullend/pkg/pagination"
)

type GigModel interface {
	WithTx(tx *sql.Tx) GigModel
	Create(clientID int64, title string, description string, budget int64, status string) (*Gig, error)
	FindByID(id int64) (*Gig, error)
	List(opts QueryOpts) (*pagination.Page[Gig], error)
	UpdateFreelancer(id int64, freelancerID int64) error
	UpdateStatus(id int64, status string) error
}

type ProposalModel interface {
	WithTx(tx *sql.Tx) ProposalModel
	Create(gigID int64, freelancerID int64, bidAmount int64, status string) (*Proposal, error)
	FindByID(id int64) (*Proposal, error)
	UpdateStatus(id int64, status string) error
}

type UserModel interface {
	WithTx(tx *sql.Tx) UserModel
	Create(email string, passwordHash string, role string, name string) (*User, error)
	FindByEmail(email string) (*User, error)
}
