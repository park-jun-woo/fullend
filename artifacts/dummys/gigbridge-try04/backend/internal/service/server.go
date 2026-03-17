package service

import (
	"github.com/gin-gonic/gin"
	"github.com/example/gigbridge/internal/middleware"
	authsvc "github.com/example/gigbridge/internal/service/auth"
	gigsvc "github.com/example/gigbridge/internal/service/gig"
	proposalsvc "github.com/example/gigbridge/internal/service/proposal"
)

// Server composes domain handlers.
type Server struct {
	Auth *authsvc.Handler
	Gig *gigsvc.Handler
	Proposal *proposalsvc.Handler
	JWTSecret string
}

// SetupRouter creates a gin.Engine that routes requests to the Server.
func SetupRouter(s *Server) *gin.Engine {
	r := gin.Default()

	// Auth group — JWT middleware extracts currentUser into context.
	auth := r.Group("/")
	auth.Use(middleware.BearerAuth(s.JWTSecret))

	auth.Handle("POST", "/gigs/:id/dispute", s.Gig.RaiseDispute)
	r.Handle("GET", "/gigs/:id", s.Gig.GetGig)
	auth.Handle("POST", "/gigs/:id/approve", s.Gig.ApproveWork)
	auth.Handle("POST", "/gigs/:id/proposals", s.Proposal.SubmitProposal)
	auth.Handle("POST", "/gigs/:id/submit-work", s.Gig.SubmitWork)
	r.Handle("POST", "/users/login", s.Auth.Login)
	auth.Handle("PUT", "/gigs/:id/publish", s.Gig.PublishGig)
	auth.Handle("POST", "/proposals/:id/accept", s.Proposal.AcceptProposal)
	r.Handle("POST", "/users/register", s.Auth.Register)
	auth.Handle("POST", "/proposals/:id/reject", s.Proposal.RejectProposal)
	r.Handle("GET", "/gigs", s.Gig.ListGigs)
	auth.Handle("POST", "/gigs", s.Gig.CreateGig)

	return r
}
