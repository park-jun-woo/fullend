package model

import (
	"context"
	"database/sql"
)

type webhookModelImpl struct {
	db *sql.DB
	tx *sql.Tx
}

func (m *webhookModelImpl) conn() interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
} {
	if m.tx != nil {
		return m.tx
	}
	return m.db
}

func NewWebhookModel(db *sql.DB) WebhookModel {
	return &webhookModelImpl{db: db}
}

func scanWebhook(row interface{ Scan(...interface{}) error }) (*Webhook, error) {
	var w Webhook
	err := row.Scan(&w.ID, &w.OrgID, &w.URL, &w.EventType, &w.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &w, nil
}

//fullend:gen ssot=db/webhooks.sql contract=baa03da
func (m *webhookModelImpl) WithTx(tx *sql.Tx) WebhookModel {
	return &webhookModelImpl{db: m.db, tx: tx}
}

//fullend:gen ssot=db/webhooks.sql contract=b8ceb42
func (m *webhookModelImpl) Create(orgID int64, url string, eventType string) (*Webhook, error) {
	row := m.conn().QueryRowContext(context.Background(),
		"INSERT INTO webhooks (org_id, url, event_type)\nVALUES ($1, $2, $3)\nRETURNING *;",
		orgID, url, eventType)
	return scanWebhook(row)
}

//fullend:gen ssot=db/webhooks.sql contract=26fb7f2
func (m *webhookModelImpl) Delete(id int64) error {
	_, err := m.conn().ExecContext(context.Background(),
		"DELETE FROM webhooks WHERE id = $1;",
		id)
	return err
}

//fullend:gen ssot=db/webhooks.sql contract=0956d71
func (m *webhookModelImpl) FindByIDAndOrgID(id int64, orgID int64) (*Webhook, error) {
	row := m.conn().QueryRowContext(context.Background(),
		"SELECT * FROM webhooks WHERE id = $1 AND org_id = $2;",
		id, orgID)
	v, err := scanWebhook(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return v, nil
}

//fullend:gen ssot=db/webhooks.sql contract=869c944
func (m *webhookModelImpl) ListByOrgID(orgID int64) ([]Webhook, error) {
	rows, err := m.conn().QueryContext(context.Background(),
		"SELECT * FROM webhooks WHERE org_id = $1;",
		orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]Webhook, 0)
	for rows.Next() {
		v, err := scanWebhook(rows)
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

//fullend:gen ssot=db/webhooks.sql contract=bafde29
func (m *webhookModelImpl) ListByOrgIDAndEventType(orgID int64, eventType string) ([]Webhook, error) {
	rows, err := m.conn().QueryContext(context.Background(),
		"SELECT * FROM webhooks WHERE org_id = $1 AND event_type = $2;",
		orgID, eventType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]Webhook, 0)
	for rows.Next() {
		v, err := scanWebhook(rows)
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
