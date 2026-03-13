package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/geul-org/zenflow/internal/model"
	"github.com/geul-org/zenflow/internal/service"
	"github.com/geul-org/fullend/pkg/authz"
	actionsvc "github.com/geul-org/zenflow/internal/service/action"
	authsvc "github.com/geul-org/zenflow/internal/service/auth"
	organizationsvc "github.com/geul-org/zenflow/internal/service/organization"
	workflowsvc "github.com/geul-org/zenflow/internal/service/workflow"
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
		Action: &actionsvc.Handler{
			DB: conn,
			ActionModel: model.NewActionModel(conn),
			WorkflowModel: model.NewWorkflowModel(conn),
		},
		Auth: &authsvc.Handler{
			DB: conn,
			UserModel: model.NewUserModel(conn),
			JWTSecret: *jwtSecret,
		},
		Organization: &organizationsvc.Handler{
			DB: conn,
			OrganizationModel: model.NewOrganizationModel(conn),
		},
		Workflow: &workflowsvc.Handler{
			DB: conn,
			ExecutionModel: model.NewExecutionModel(conn),
			UserModel: model.NewUserModel(conn),
			WorkflowModel: model.NewWorkflowModel(conn),
		},
		JWTSecret: *jwtSecret,
	}

	r := service.SetupRouter(server)
	log.Printf("server listening on %s", *addr)
	log.Fatal(r.Run(*addr))
}
