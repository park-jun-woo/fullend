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
	err := row.Scan(&e.ID, &e.WorkflowID, &e.OrgID, &e.Status, &e.CreditsSpent, &e.ReportKey, &e.ExecutedAt)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

//fullend:gen ssot=db/execution_logs.sql contract=03852b3
func (m *executionlogModelImpl) WithTx(tx *sql.Tx) ExecutionLogModel {
	return &executionlogModelImpl{db: m.db, tx: tx}
}

//fullend:gen ssot=db/execution_logs.sql contract=3aa3a8a
func (m *executionlogModelImpl) Create(workflowID int64, orgID int64, status string, creditsSpent int64) (*ExecutionLog, error) {
	row := m.conn().QueryRowContext(context.Background(),
		"INSERT INTO execution_logs (workflow_id, org_id, status, credits_spent)\nVALUES ($1, $2, $3, $4)\nRETURNING *;",
		workflowID, orgID, status, creditsSpent)
	return scanExecutionLog(row)
}

//fullend:gen ssot=db/execution_logs.sql contract=8386d35
func (m *executionlogModelImpl) CreateWithReport(workflowID int64, orgID int64, status string, creditsSpent int64, reportKey string) (*ExecutionLog, error) {
	row := m.conn().QueryRowContext(context.Background(),
		"INSERT INTO execution_logs (workflow_id, org_id, status, credits_spent, report_key)\nVALUES ($1, $2, $3, $4, $5)\nRETURNING *;",
		workflowID, orgID, status, creditsSpent, reportKey)
	return scanExecutionLog(row)
}

//fullend:gen ssot=db/execution_logs.sql contract=797a98e
func (m *executionlogModelImpl) FindByIDAndOrgID(id int64, orgID int64) (*ExecutionLog, error) {
	row := m.conn().QueryRowContext(context.Background(),
		"SELECT * FROM execution_logs WHERE id = $1 AND org_id = $2;",
		id, orgID)
	v, err := scanExecutionLog(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return v, nil
}

//fullend:gen ssot=db/execution_logs.sql contract=2cb689e
func (m *executionlogModelImpl) ListByWorkflowID(workflowID int64) ([]ExecutionLog, error) {
	rows, err := m.conn().QueryContext(context.Background(),
		"SELECT * FROM execution_logs WHERE workflow_id = $1 ORDER BY executed_at DESC;",
		workflowID)
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
