package service

import (
	"github.com/gin-gonic/gin"
	"github.com/example/zenflow/internal/middleware"
	authsvc "github.com/example/zenflow/internal/service/auth"
	logsvc "github.com/example/zenflow/internal/service/log"
	workflowsvc "github.com/example/zenflow/internal/service/workflow"
)

// Server composes domain handlers.
type Server struct {
	Auth *authsvc.Handler
	Log *logsvc.Handler
	Workflow *workflowsvc.Handler
	JWTSecret string
}

// SetupRouter creates a gin.Engine that routes requests to the Server.
func SetupRouter(s *Server) *gin.Engine {
	r := gin.Default()

	// Auth group — JWT middleware extracts currentUser into context.
	auth := r.Group("/")
	auth.Use(middleware.BearerAuth(s.JWTSecret))

	auth.POST("/workflows/:id/pause", s.Workflow.PauseWorkflow)
	auth.GET("/workflows/:id", s.Workflow.GetWorkflow)
	auth.POST("/workflows/:id/actions", s.Workflow.AddAction)
	auth.POST("/workflows/:id/execute", s.Workflow.ExecuteWorkflow)
	r.POST("/auth/register", s.Auth.Register)
	auth.GET("/execution-logs", s.Log.ListExecutionLogs)
	auth.GET("/workflows", s.Workflow.ListWorkflows)
	auth.POST("/workflows", s.Workflow.CreateWorkflow)
	auth.POST("/workflows/:id/activate", s.Workflow.ActivateWorkflow)
	r.POST("/auth/login", s.Auth.Login)
	auth.POST("/workflows/:id/archive", s.Workflow.ArchiveWorkflow)

	return r
}
