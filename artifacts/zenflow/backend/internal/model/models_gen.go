package model

import (
	"database/sql"
)

type ActionModel interface {
	WithTx(tx *sql.Tx) ActionModel
	Create(workflowID int64, actionType string, payloadTemplate string, sequenceOrder int64) (*Action, error)
	ListByWorkflow(workflowID int64) ([]Action, error)
}

type ExecutionModel interface {
	WithTx(tx *sql.Tx) ExecutionModel
	Create(workflowID int64, orgID int64, status string, creditsSpent int64) (*Execution, error)
}

type OrganizationModel interface {
	WithTx(tx *sql.Tx) OrganizationModel
	Create(name string, planType string, creditsBalance int64) (*Organization, error)
	FindByID(id int64) (*Organization, error)
	UpdateCredits(id int64, creditsBalance int64) error
}

type UserModel interface {
	WithTx(tx *sql.Tx) UserModel
	Create(orgID int64, email string, passwordHash string, role string) (*User, error)
	FindByEmail(email string) (*User, error)
	FindByID(id int64) (*User, error)
}

type WorkflowModel interface {
	WithTx(tx *sql.Tx) WorkflowModel
	Create(orgID int64, title string, triggerEvent string) (*Workflow, error)
	FindByID(id int64) (*Workflow, error)
	FindByIDAndOrg(id int64, orgID int64) (*Workflow, error)
	List(orgID int64) ([]Workflow, error)
	UpdateStatus(id int64, status string) error
}
