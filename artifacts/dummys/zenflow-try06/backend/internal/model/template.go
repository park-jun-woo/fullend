package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/park-jun-woo/fullend/pkg/pagination"
)

type templateModelImpl struct {
	db *sql.DB
	tx *sql.Tx
}

func (m *templateModelImpl) conn() interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
} {
	if m.tx != nil {
		return m.tx
	}
	return m.db
}

func NewTemplateModel(db *sql.DB) TemplateModel {
	return &templateModelImpl{db: db}
}

func scanTemplate(row interface{ Scan(...interface{}) error }) (*Template, error) {
	var t Template
	err := row.Scan(&t.ID, &t.SourceWorkflowID, &t.OrgID, &t.Title, &t.Description, &t.Category, &t.CloneCount, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

//fullend:gen ssot=db/templates.sql contract=6124b76
func (m *templateModelImpl) WithTx(tx *sql.Tx) TemplateModel {
	return &templateModelImpl{db: m.db, tx: tx}
}

//fullend:gen ssot=db/templates.sql contract=623e156
func (m *templateModelImpl) Create(sourceWorkflowID int64, orgID int64, title string, description string, category string) (*Template, error) {
	row := m.conn().QueryRowContext(context.Background(),
		"INSERT INTO templates (source_workflow_id, org_id, title, description, category)\nVALUES ($1, $2, $3, $4, $5)\nRETURNING *;",
		sourceWorkflowID, orgID, title, description, category)
	return scanTemplate(row)
}

//fullend:gen ssot=db/templates.sql contract=1c6a35a
func (m *templateModelImpl) FindByID(id int64) (*Template, error) {
	row := m.conn().QueryRowContext(context.Background(),
		"SELECT * FROM templates WHERE id = $1;",
		id)
	v, err := scanTemplate(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return v, nil
}

//fullend:gen ssot=db/templates.sql contract=6c1ccc2
func (m *templateModelImpl) FindBySourceWorkflowID(sourceWorkflowID int64) (*Template, error) {
	row := m.conn().QueryRowContext(context.Background(),
		"SELECT * FROM templates WHERE source_workflow_id = $1;",
		sourceWorkflowID)
	v, err := scanTemplate(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return v, nil
}

//fullend:gen ssot=db/templates.sql contract=c64bc6d
func (m *templateModelImpl) IncrementCloneCount(id int64) error {
	_, err := m.conn().ExecContext(context.Background(),
		"UPDATE templates SET clone_count = clone_count + 1 WHERE id = $1;",
		id)
	return err
}

//fullend:gen ssot=db/templates.sql contract=077717c
func (m *templateModelImpl) List(opts QueryOpts) (*pagination.Cursor[Template], error) {
	requestedLimit := opts.Limit
	opts.Limit = requestedLimit + 1

	selectSQL, selectArgs := BuildSelectQuery("templates", "", 0, opts)
	rows, err := m.conn().QueryContext(context.Background(), selectSQL, selectArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]Template, 0)
	for rows.Next() {
		v, err := scanTemplate(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *v)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	hasNext := len(items) > requestedLimit
	var nextCursor string
	if hasNext {
		items = items[:requestedLimit]
	}
	if len(items) > 0 {
		nextCursor = fmt.Sprintf("%v", items[len(items)-1].ID)
	}
	return &pagination.Cursor[Template]{Items: items, NextCursor: nextCursor, HasNext: hasNext}, nil
}
