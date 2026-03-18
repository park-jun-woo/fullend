package service

import (
	"github.com/gin-gonic/gin"
	"github.com/example/zenflow/internal/middleware"
	actionsvc "github.com/example/zenflow/internal/service/action"
	authsvc "github.com/example/zenflow/internal/service/auth"
	logsvc "github.com/example/zenflow/internal/service/log"
	organizationsvc "github.com/example/zenflow/internal/service/organization"
	schedulesvc "github.com/example/zenflow/internal/service/schedule"
	templatesvc "github.com/example/zenflow/internal/service/template"
	webhooksvc "github.com/example/zenflow/internal/service/webhook"
	workflowsvc "github.com/example/zenflow/internal/service/workflow"
)

// Server composes domain handlers.
type Server struct {
	Action *actionsvc.Handler
	Auth *authsvc.Handler
	Log *logsvc.Handler
	Organization *organizationsvc.Handler
	Schedule *schedulesvc.Handler
	Template *templatesvc.Handler
	Webhook *webhooksvc.Handler
	Workflow *workflowsvc.Handler
	JWTSecret string
}

// SetupRouter creates a gin.Engine that routes requests to the Server.
func SetupRouter(s *Server) *gin.Engine {
	r := gin.Default()

	// Auth group — JWT middleware extracts currentUser into context.
	auth := r.Group("/")
	auth.Use(middleware.BearerAuth(s.JWTSecret))

	r.Handle("GET", "/templates", s.Template.ListTemplates)
	auth.Handle("POST", "/templates", s.Template.PublishTemplate)
	auth.Handle("POST", "/workflows/:id/actions", s.Action.AddAction)
	auth.Handle("GET", "/workflows/:id", s.Workflow.GetWorkflow)
	r.Handle("POST", "/users/login", s.Auth.Login)
	auth.Handle("POST", "/workflows/:id/activate", s.Workflow.ActivateWorkflow)
	auth.Handle("POST", "/workflows/:id/new-version", s.Workflow.CreateWorkflowVersion)
	r.Handle("GET", "/templates/:id", s.Template.GetTemplate)
	auth.Handle("POST", "/workflows/:id/archive", s.Workflow.ArchiveWorkflow)
	auth.Handle("GET", "/workflows/:id/versions", s.Workflow.ListWorkflowVersions)
	auth.Handle("POST", "/workflows/:id/execute-with-report", s.Workflow.ExecuteWithReport)
	r.Handle("POST", "/organizations", s.Organization.CreateOrganization)
	r.Handle("POST", "/users/register", s.Auth.Register)
	auth.Handle("GET", "/execution-logs/:id/report", s.Log.GetExecutionReport)
	auth.Handle("GET", "/webhooks", s.Webhook.ListWebhooks)
	auth.Handle("POST", "/webhooks", s.Webhook.CreateWebhook)
	auth.Handle("DELETE", "/webhooks/:id", s.Webhook.DeleteWebhook)
	auth.Handle("POST", "/workflows/:id/pause", s.Workflow.PauseWorkflow)
	auth.Handle("POST", "/workflows/:id/execute", s.Workflow.ExecuteWorkflow)
	auth.Handle("GET", "/workflows", s.Workflow.ListWorkflows)
	auth.Handle("POST", "/workflows", s.Workflow.CreateWorkflow)
	auth.Handle("POST", "/templates/:id/clone", s.Template.CloneTemplate)
	auth.Handle("GET", "/workflows/:id/logs", s.Log.ListExecutionLogs)
	auth.Handle("DELETE", "/workflows/:id/schedule", s.Schedule.DeleteSchedule)
	auth.Handle("GET", "/workflows/:id/schedule", s.Schedule.GetSchedule)
	auth.Handle("POST", "/workflows/:id/schedule", s.Schedule.SetSchedule)

	return r
}
