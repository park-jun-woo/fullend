package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/example/gigbridge/internal/model"
	"github.com/example/gigbridge/internal/service"
	"github.com/park-jun-woo/fullend/pkg/authz"
	authsvc "github.com/example/gigbridge/internal/service/auth"
	gigsvc "github.com/example/gigbridge/internal/service/gig"
	proposalsvc "github.com/example/gigbridge/internal/service/proposal"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	dsn := flag.String("dsn", "postgres://localhost:5432/app?sslmode=disable", "database connection string")
	dbDriver := flag.String("db", "postgres", "database driver (postgres, mysql)")
	jwtSecretDefault := os.Getenv("JWT_SECRET")
	if jwtSecretDefault == "" {
		jwtSecretDefault = "secret"
	}
	jwtSecret := flag.String("jwt-secret", jwtSecretDefault, "JWT signing secret")
	flag.Parse()

	conn, err := sql.Open(*dbDriver, *dsn)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		log.Fatalf("database ping failed: %v", err)
	}

	os.Setenv("JWT_SECRET", *jwtSecret)

	if err := authz.Init(conn, []authz.OwnershipMapping{
		{Resource: "gig", Table: "gigs", Column: "client_id"},
		{Resource: "gig_assignee", Table: "gigs", Column: "freelancer_id"},
		{Resource: "proposal", Table: "proposals", Column: "freelancer_id"},
	}); err != nil {
		log.Fatalf("authz init failed: %v", err)
	}

	server := &service.Server{
		Auth: &authsvc.Handler{
			DB: conn,
			UserModel: model.NewUserModel(conn),
			JWTSecret: *jwtSecret,
		},
		Gig: &gigsvc.Handler{
			DB: conn,
			GigModel: model.NewGigModel(conn),
		},
		Proposal: &proposalsvc.Handler{
			DB: conn,
			GigModel: model.NewGigModel(conn),
			ProposalModel: model.NewProposalModel(conn),
		},
		JWTSecret: *jwtSecret,
	}

	r := service.SetupRouter(server)
	log.Printf("server listening on %s", *addr)
	log.Fatal(r.Run(*addr))
}
