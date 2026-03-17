package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/example/zenflow/internal/model"
	"github.com/example/zenflow/internal/service"
	"github.com/park-jun-woo/fullend/pkg/authz"
	authsvc "github.com/example/zenflow/internal/service/auth"
	logsvc "github.com/example/zenflow/internal/service/log"
	workflowsvc "github.com/example/zenflow/internal/service/workflow"
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
		{Resource: "workflow", Table: "workflows", Column: "org_id"},
	}); err != nil {
		log.Fatalf("authz init failed: %v", err)
	}

	server := &service.Server{
		Auth: &authsvc.Handler{
			DB: conn,
			OrganizationModel: model.NewOrganizationModel(conn),
			UserModel: model.NewUserModel(conn),
			JWTSecret: *jwtSecret,
		},
		Log: &logsvc.Handler{
			ExecutionLogModel: model.NewExecutionLogModel(conn),
		},
		Workflow: &workflowsvc.Handler{
			DB: conn,
			ActionModel: model.NewActionModel(conn),
			ExecutionLogModel: model.NewExecutionLogModel(conn),
			OrganizationModel: model.NewOrganizationModel(conn),
			WorkflowModel: model.NewWorkflowModel(conn),
		},
		JWTSecret: *jwtSecret,
	}

	r := service.SetupRouter(server)
	log.Printf("server listening on %s", *addr)
	log.Fatal(r.Run(*addr))
}
