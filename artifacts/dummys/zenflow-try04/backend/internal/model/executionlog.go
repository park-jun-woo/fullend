package model

import (
	"context"
	"database/sql"
)

type executionlogModelImpl struct {
	db *sql.DB
	tx *sql.Tx
}

func (m *executionlogModelImpl) conn() interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
} {
	if m.tx != nil {
		return m.tx
	}
	return m.db
}

func NewExecutionLogModel(db *sql.DB) ExecutionLogModel {
	return &executionlogModelImpl{db: db}
}

func scanExecutionLog(row interface{ Scan(...interface{}) error }) (*ExecutionLog, error) {
	var e ExecutionLog
	err := row.Scan(&e.ID, &e.WorkflowID, &e.OrgID, &e.Status, &e.CreditsSpent, &e.ExecutedAt)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

//fullend:gen ssot=db/execution_logs.sql contract=03852b3
func (m *executionlogModelImpl) WithTx(tx *sql.Tx) ExecutionLogModel {
	return &executionlogModelImpl{db: m.db, tx: tx}
}

//fullend:gen ssot=db/execution_logs.sql contract=073d800
func (m *executionlogModelImpl) Create(workflowID int64, orgID int64, status string) (*ExecutionLog, error) {
	row := m.conn().QueryRowContext(context.Background(),
		"INSERT INTO execution_logs (workflow_id, org_id, status)\nVALUES ($1, $2, $3)\nRETURNING *;",
		workflowID, orgID, status)
	return scanExecutionLog(row)
}
