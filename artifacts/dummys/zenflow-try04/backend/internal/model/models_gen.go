package model

import (
	"database/sql"
)

type ActionModel interface {
	WithTx(tx *sql.Tx) ActionModel
	Create(workflowID int64, actionType string, payloadTemplate string, sequenceOrder int64) (*Action, error)
	ListByWorkflowID(workflowID int64) ([]Action, error)
}

type ExecutionLogModel interface {
	WithTx(tx *sql.Tx) ExecutionLogModel
	Create(workflowID int64, orgID int64, status string) (*ExecutionLog, error)
}

type OrganizationModel interface {
	WithTx(tx *sql.Tx) OrganizationModel
	DeductCredit(id int64) error
	FindByID(id int64) (*Organization, error)
}

type UserModel interface {
	WithTx(tx *sql.Tx) UserModel
	Create(email string, passwordHash string, orgID int64, role string) (*User, error)
	FindByEmail(email string) (*User, error)
	FindByID(id int64) (*User, error)
}

type WorkflowModel interface {
	WithTx(tx *sql.Tx) WorkflowModel
	Create(orgID int64, title string, triggerEvent string) (*Workflow, error)
	FindByID(id int64) (*Workflow, error)
	ListByOrgID(orgID int64) ([]Workflow, error)
	UpdateStatus(id int64, status string) error
}
