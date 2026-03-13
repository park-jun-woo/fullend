package service

import (
	"github.com/gin-gonic/gin"
	"github.com/geul-org/zenflow/internal/middleware"
	actionsvc "github.com/geul-org/zenflow/internal/service/action"
	authsvc "github.com/geul-org/zenflow/internal/service/auth"
	organizationsvc "github.com/geul-org/zenflow/internal/service/organization"
	workflowsvc "github.com/geul-org/zenflow/internal/service/workflow"
)

// Server composes domain handlers.
type Server struct {
	Action *actionsvc.Handler
	Auth *authsvc.Handler
	Organization *organizationsvc.Handler
	Workflow *workflowsvc.Handler
	JWTSecret string
}

// SetupRouter creates a gin.Engine that routes requests to the Server.
func SetupRouter(s *Server) *gin.Engine {
	r := gin.Default()

	// Auth group — JWT middleware extracts currentUser into context.
	auth := r.Group("/")
	auth.Use(middleware.BearerAuth(s.JWTSecret))

	auth.POST("/workflows/:id/activate", s.Workflow.ActivateWorkflow)
	r.POST("/auth/register", s.Auth.Register)
	auth.POST("/workflows/:id/archive", s.Workflow.ArchiveWorkflow)
	auth.POST("/workflows/:id/execute", s.Workflow.ExecuteWorkflow)
	auth.POST("/workflows/:id/pause", s.Workflow.PauseWorkflow)
	r.POST("/auth/login", s.Auth.Login)
	r.POST("/organizations", s.Organization.CreateOrganization)
	auth.GET("/workflows/:id/actions", s.Action.ListActions)
	auth.POST("/workflows/:id/actions", s.Action.CreateAction)
	auth.GET("/workflows", s.Workflow.ListWorkflows)
	auth.POST("/workflows", s.Workflow.CreateWorkflow)
	auth.GET("/workflows/:id", s.Workflow.GetWorkflow)

	return r
}
