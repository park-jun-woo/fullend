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

//fullend:gen ssot=db/execution_logs.sql contract=c4344d1
func (m *executionlogModelImpl) Create(workflowID int64, orgID int64) (*ExecutionLog, error) {
	row := m.conn().QueryRowContext(context.Background(),
		"INSERT INTO execution_logs (workflow_id, org_id)\nVALUES ($1, $2)\nRETURNING *;",
		workflowID, orgID)
	return scanExecutionLog(row)
}

//fullend:gen ssot=db/execution_logs.sql contract=73c7f86
func (m *executionlogModelImpl) ListByOrg(orgID int64) ([]ExecutionLog, error) {
	rows, err := m.conn().QueryContext(context.Background(),
		"SELECT * FROM execution_logs WHERE org_id = $1 ORDER BY executed_at DESC;",
		orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]ExecutionLog, 0)
	for rows.Next() {
		v, err := scanExecutionLog(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *v)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
