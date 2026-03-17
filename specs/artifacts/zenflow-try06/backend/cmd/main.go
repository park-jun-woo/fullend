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
	actionsvc "github.com/example/zenflow/internal/service/action"
	authsvc "github.com/example/zenflow/internal/service/auth"
	logsvc "github.com/example/zenflow/internal/service/log"
	organizationsvc "github.com/example/zenflow/internal/service/organization"
	schedulesvc "github.com/example/zenflow/internal/service/schedule"
	templatesvc "github.com/example/zenflow/internal/service/template"
	webhooksvc "github.com/example/zenflow/internal/service/webhook"
	workflowsvc "github.com/example/zenflow/internal/service/workflow"
	"context"
	"encoding/json"
	"github.com/park-jun-woo/fullend/pkg/queue"
	"fmt"
	"github.com/park-jun-woo/fullend/pkg/session"
	"github.com/park-jun-woo/fullend/pkg/file"
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

	if err := queue.Init(context.Background(), "postgres", conn); err != nil {
		log.Fatalf("queue init failed: %v", err)
	}
	defer queue.Close()

	sm, err := session.NewPostgresSession(context.Background(), conn)
	if err != nil {
		log.Fatalf("session init failed: %v", err)
	}
	session.Init(sm)
	file.Init(file.NewLocalFile("./uploads"))
	server := &service.Server{
		Action: &actionsvc.Handler{
			DB: conn,
			ActionModel: model.NewActionModel(conn),
			WorkflowModel: model.NewWorkflowModel(conn),
		},
		Auth: &authsvc.Handler{
			DB: conn,
			OrganizationModel: model.NewOrganizationModel(conn),
			UserModel: model.NewUserModel(conn),
			JWTSecret: *jwtSecret,
		},
		Log: &logsvc.Handler{
			ExecutionLogModel: model.NewExecutionLogModel(conn),
			WorkflowModel: model.NewWorkflowModel(conn),
		},
		Organization: &organizationsvc.Handler{
			DB: conn,
			OrganizationModel: model.NewOrganizationModel(conn),
		},
		Schedule: &schedulesvc.Handler{
			WorkflowModel: model.NewWorkflowModel(conn),
		},
		Template: &templatesvc.Handler{
			DB: conn,
			ActionModel: model.NewActionModel(conn),
			OrganizationModel: model.NewOrganizationModel(conn),
			TemplateModel: model.NewTemplateModel(conn),
			WorkflowModel: model.NewWorkflowModel(conn),
		},
		Webhook: &webhooksvc.Handler{
			DB: conn,
			WebhookModel: model.NewWebhookModel(conn),
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

	queue.Subscribe("workflow.executed", func(ctx context.Context, msg []byte) error {
		var message webhooksvc.WorkflowExecutedMessage
		if err := json.Unmarshal(msg, &message); err != nil {
			return fmt.Errorf("unmarshal: %w", err)
		}
		return server.Webhook.OnWorkflowExecuted(ctx, message)
	})

	go queue.Start(context.Background())

	r := service.SetupRouter(server)
	log.Printf("server listening on %s", *addr)
	log.Fatal(r.Run(*addr))
}
