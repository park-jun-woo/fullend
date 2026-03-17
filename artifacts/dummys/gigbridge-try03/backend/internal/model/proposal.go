package model

import (
	"context"
	"database/sql"
)

type proposalModelImpl struct {
	db *sql.DB
	tx *sql.Tx
}

func (m *proposalModelImpl) conn() interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
} {
	if m.tx != nil {
		return m.tx
	}
	return m.db
}

func NewProposalModel(db *sql.DB) ProposalModel {
	return &proposalModelImpl{db: db}
}

func scanProposal(row interface{ Scan(...interface{}) error }) (*Proposal, error) {
	var p Proposal
	err := row.Scan(&p.ID, &p.GigID, &p.FreelancerID, &p.BidAmount, &p.Status)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

//fullend:gen ssot=db/proposals.sql contract=6c2a2cd
func (m *proposalModelImpl) WithTx(tx *sql.Tx) ProposalModel {
	return &proposalModelImpl{db: m.db, tx: tx}
}

//fullend:gen ssot=db/proposals.sql contract=6a98c16
func (m *proposalModelImpl) Create(gigID int64, freelancerID int64, bidAmount int64, status string) (*Proposal, error) {
	row := m.conn().QueryRowContext(context.Background(),
		"INSERT INTO proposals (gig_id, freelancer_id, bid_amount, status)\nVALUES ($1, $2, $3, $4)\nRETURNING id, gig_id, freelancer_id, bid_amount, status;",
		gigID, freelancerID, bidAmount, status)
	return scanProposal(row)
}

//fullend:gen ssot=db/proposals.sql contract=cf2a010
func (m *proposalModelImpl) FindByID(id int64) (*Proposal, error) {
	row := m.conn().QueryRowContext(context.Background(),
		"SELECT id, gig_id, freelancer_id, bid_amount, status FROM proposals WHERE id = $1;",
		id)
	v, err := scanProposal(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return v, nil
}

//fullend:gen ssot=db/proposals.sql contract=2b9704e
func (m *proposalModelImpl) UpdateStatus(id int64, status string) error {
	_, err := m.conn().ExecContext(context.Background(),
		"UPDATE proposals SET status = $2 WHERE id = $1;",
		id, status)
	return err
}
