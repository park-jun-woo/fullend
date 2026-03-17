package model

import (
	"time"
)

type Action struct {
	ID           int64 `json:"id"`
	WorkflowID   int64 `json:"workflow_id"`
	ActionType   string `json:"action_type"`
	PayloadTemplate string `json:"payload_template"`
	SequenceOrder int64 `json:"sequence_order"`
}

type ExecutionLog struct {
	ID           int64 `json:"id"`
	WorkflowID   int64 `json:"workflow_id"`
	OrgID        int64 `json:"org_id"`
	Status       string `json:"status"`
	CreditsSpent int64 `json:"credits_spent"`
	ReportKey    string `json:"report_key"`
	ExecutedAt   time.Time `json:"executed_at"`
}

type Organization struct {
	ID           int64 `json:"id"`
	Name         string `json:"name"`
	PlanType     string `json:"plan_type"`
	CreditsBalance int64 `json:"credits_balance"`
}

type Template struct {
	ID           int64 `json:"id"`
	SourceWorkflowID int64 `json:"source_workflow_id"`
	OrgID        int64 `json:"org_id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Category     string `json:"category"`
	CloneCount   int64 `json:"clone_count"`
	CreatedAt    time.Time `json:"created_at"`
}

type User struct {
	ID           int64 `json:"id"`
	OrgID        int64 `json:"org_id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	Role         string `json:"role"`
	Name         string `json:"name"`
}

type Webhook struct {
	ID           int64 `json:"id"`
	OrgID        int64 `json:"org_id"`
	URL          string `json:"url"`
	EventType    string `json:"event_type"`
	CreatedAt    time.Time `json:"created_at"`
}

type Workflow struct {
	ID           int64 `json:"id"`
	OrgID        int64 `json:"org_id"`
	Title        string `json:"title"`
	TriggerEvent string `json:"trigger_event"`
	Status       string `json:"status"`
	Version      int64 `json:"version"`
	RootWorkflowID int64 `json:"root_workflow_id"`
	CreatedAt    time.Time `json:"created_at"`
}
