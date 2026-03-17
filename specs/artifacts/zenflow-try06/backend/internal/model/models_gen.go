package model

import (
	"database/sql"

	"github.com/park-jun-woo/fullend/pkg/pagination"
)

type ActionModel interface {
	WithTx(tx *sql.Tx) ActionModel
	CopyToWorkflow(workflowID int64, sourceWorkflowID int64) error
	Create(workflowID int64, actionType string, payloadTemplate string, sequenceOrder int64) (*Action, error)
	ListByWorkflowID(workflowID int64) ([]Action, error)
}

type ExecutionLogModel interface {
	WithTx(tx *sql.Tx) ExecutionLogModel
	Create(workflowID int64, orgID int64, status string, creditsSpent int64) (*ExecutionLog, error)
	CreateWithReport(workflowID int64, orgID int64, status string, creditsSpent int64, reportKey string) (*ExecutionLog, error)
	FindByIDAndOrgID(id int64, orgID int64) (*ExecutionLog, error)
	ListByWorkflowID(workflowID int64) ([]ExecutionLog, error)
}

type OrganizationModel interface {
	WithTx(tx *sql.Tx) OrganizationModel
	Create(name string, planType string, creditsBalance int64) (*Organization, error)
	FindByID(id int64) (*Organization, error)
	UpdateCredits(id int64, creditsBalance int64) error
}

type TemplateModel interface {
	WithTx(tx *sql.Tx) TemplateModel
	Create(sourceWorkflowID int64, orgID int64, title string, description string, category string) (*Template, error)
	FindByID(id int64) (*Template, error)
	FindBySourceWorkflowID(sourceWorkflowID int64) (*Template, error)
	IncrementCloneCount(id int64) error
	List(opts QueryOpts) (*pagination.Cursor[Template], error)
}

type UserModel interface {
	WithTx(tx *sql.Tx) UserModel
	Create(orgID int64, email string, passwordHash string, role string, name string) (*User, error)
	FindByEmail(email string) (*User, error)
}

type WebhookModel interface {
	WithTx(tx *sql.Tx) WebhookModel
	Create(orgID int64, url string, eventType string) (*Webhook, error)
	Delete(id int64) error
	FindByIDAndOrgID(id int64, orgID int64) (*Webhook, error)
	ListByOrgID(orgID int64) ([]Webhook, error)
	ListByOrgIDAndEventType(orgID int64, eventType string) ([]Webhook, error)
}

type WorkflowModel interface {
	WithTx(tx *sql.Tx) WorkflowModel
	Create(orgID int64, title string, triggerEvent string, status string) (*Workflow, error)
	CreateVersion(orgID int64, title string, triggerEvent string, status string, version int64, rootWorkflowID int64) (*Workflow, error)
	FindByID(id int64) (*Workflow, error)
	FindByIDAndOrgID(id int64, orgID int64) (*Workflow, error)
	ListByOrgID(orgID int64) ([]Workflow, error)
	ListVersions(rootWorkflowID int64, orgID int64) ([]Workflow, error)
	UpdateStatus(id int64, status string) error
}
