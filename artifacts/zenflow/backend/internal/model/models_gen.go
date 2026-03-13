package model

import (
	"database/sql"
	"encoding/json"
)

type ActionModel interface {
	WithTx(tx *sql.Tx) ActionModel
	Create(workflowID int64, actionType string, payloadTemplate json.RawMessage, sequenceOrder int64) (*Action, error)
	ListByWorkflowID(workflowID int64) ([]Action, error)
}

type ExecutionModel interface {
	WithTx(tx *sql.Tx) ExecutionModel
	Create(workflowID int64, orgID int64, status string, creditsSpent int64) (*Execution, error)
}

type OrganizationModel interface {
	WithTx(tx *sql.Tx) OrganizationModel
	Create(name string, planType string, creditsBalance int64) (*Organization, error)
}

type UserModel interface {
	WithTx(tx *sql.Tx) UserModel
	Create(orgID int64, email string, passwordHash string, role string) (*User, error)
	FindByEmail(email string) (*User, error)
	FindByID(id int64) (*User, error)
}

type WorkflowModel interface {
	WithTx(tx *sql.Tx) WorkflowModel
	Create(orgID int64, title string, triggerEvent string, status string) (*Workflow, error)
	FindByID(id int64) (*Workflow, error)
	ListByOrgID(orgID int64) ([]Workflow, error)
	UpdateStatus(status string, id int64) error
}
