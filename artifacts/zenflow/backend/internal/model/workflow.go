package model

import (
	"context"
	"database/sql"
)

type workflowModelImpl struct {
	db *sql.DB
	tx *sql.Tx
}

func (m *workflowModelImpl) conn() interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
} {
	if m.tx != nil {
		return m.tx
	}
	return m.db
}

func NewWorkflowModel(db *sql.DB) WorkflowModel {
	return &workflowModelImpl{db: db}
}

func scanWorkflow(row interface{ Scan(...interface{}) error }) (*Workflow, error) {
	var w Workflow
	err := row.Scan(&w.ID, &w.OrgID, &w.Title, &w.TriggerEvent, &w.Status, &w.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &w, nil
}

//fullend:gen ssot=db/workflows.sql contract=94ac658
func (m *workflowModelImpl) WithTx(tx *sql.Tx) WorkflowModel {
	return &workflowModelImpl{db: m.db, tx: tx}
}

//fullend:gen ssot=db/workflows.sql contract=0c3ed66
func (m *workflowModelImpl) Create(orgID int64, title string, triggerEvent string, status string) (*Workflow, error) {
	row := m.conn().QueryRowContext(context.Background(),
		"INSERT INTO workflows (org_id, title, trigger_event, status)\nVALUES ($1, $2, $3, $4)\nRETURNING *;",
		orgID, title, triggerEvent, status)
	return scanWorkflow(row)
}

//fullend:gen ssot=db/workflows.sql contract=8e417ed
func (m *workflowModelImpl) FindByID(id int64) (*Workflow, error) {
	row := m.conn().QueryRowContext(context.Background(),
		"SELECT * FROM workflows WHERE id = $1;",
		id)
	v, err := scanWorkflow(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return v, nil
}

//fullend:gen ssot=db/workflows.sql contract=d4de074
func (m *workflowModelImpl) ListByOrgID(orgID int64) ([]Workflow, error) {
	rows, err := m.conn().QueryContext(context.Background(),
		"SELECT * FROM workflows WHERE org_id = $1;",
		orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]Workflow, 0)
	for rows.Next() {
		v, err := scanWorkflow(rows)
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

//fullend:gen ssot=db/workflows.sql contract=dddd157
func (m *workflowModelImpl) UpdateStatus(status string, id int64) error {
	_, err := m.conn().ExecContext(context.Background(),
		"UPDATE workflows SET status = $1 WHERE id = $2;",
		status, id)
	return err
}
