//ff:func feature=gen-gogin type=generator
//ff:what creates cmd/main.go with domain handler initialization

package gogin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/geul-org/fullend/internal/policy"
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// generateMainWithDomains creates cmd/main.go with domain handler initialization.
func generateMainWithDomains(artifactsDir string, serviceFuncs []ssacparser.ServiceFunc, modulePath string, queueBackend string, policies []*policy.Policy) error {
	if modulePath == "" {
		base := filepath.Base(artifactsDir)
		modulePath = base + "/backend"
	}

	goModPath := filepath.Join(artifactsDir, "backend", "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Join(artifactsDir, "backend"), 0755); err != nil {
			return err
		}
		goModContent := fmt.Sprintf("module %s\n\ngo 1.22\n\nrequire github.com/gin-gonic/gin v1.10.0\n", modulePath)
		if err := os.WriteFile(goModPath, []byte(goModContent), 0644); err != nil {
			return err
		}
	}

	if err := os.MkdirAll(filepath.Join(artifactsDir, "backend", "cmd"), 0755); err != nil {
		return err
	}

	domains := uniqueDomains(serviceFuncs)
	flatModels := collectModelsForDomain(serviceFuncs, "")
	anyNeedsAuth := false
	for _, d := range domains {
		if domainNeedsAuth(serviceFuncs, d) {
			anyNeedsAuth = true
			break
		}
	}

	// Build init block.
	var initLines []string

	// Flat model fields.
	for _, m := range flatModels {
		fieldName := ucFirst(lcFirst(m) + "Model")
		initLines = append(initLines, fmt.Sprintf("\t\t%s: model.New%sModel(conn),", fieldName, m))
	}

	// Domain handler fields.
	for _, domain := range domains {
		domainModels := collectModelsForDomain(serviceFuncs, domain)
		fieldName := ucFirst(domain)

		var handlerLines []string
		if domainNeedsDB(serviceFuncs, domain) {
			handlerLines = append(handlerLines, "\t\t\tDB: conn,")
		}
		for _, m := range domainModels {
			mFieldName := ucFirst(lcFirst(m) + "Model")
			handlerLines = append(handlerLines, fmt.Sprintf("\t\t\t%s: model.New%sModel(conn),", mFieldName, m))
		}
		if domainNeedsJWTSecret(serviceFuncs, domain) {
			handlerLines = append(handlerLines, "\t\t\tJWTSecret: *jwtSecret,")
		}
		initLines = append(initLines, fmt.Sprintf("\t\t%s: &%ssvc.Handler{", fieldName, domain))
		initLines = append(initLines, handlerLines...)
		initLines = append(initLines, "\t\t},")
	}

	// Add JWTSecret to Server struct if bearer auth is used.
	if anyNeedsAuth {
		initLines = append(initLines, "\t\tJWTSecret: *jwtSecret,")
	}

	initBlock := strings.Join(initLines, "\n")
	if initBlock == "" {
		initBlock = "\t\t// No models detected"
	}

	// Build domain imports.
	var extraImports []string
	extraImports = append(extraImports, fmt.Sprintf("\n\t\"%s/internal/model\"", modulePath))
	extraImports = append(extraImports, fmt.Sprintf("\t\"%s/internal/service\"", modulePath))
	if anyNeedsAuth {
		extraImports = append(extraImports, "\t\"github.com/geul-org/fullend/pkg/authz\"")
	}
	for _, d := range domains {
		extraImports = append(extraImports, fmt.Sprintf("\t%ssvc \"%s/internal/service/%s\"", d, modulePath, d))
	}
	importBlock := strings.Join(extraImports, "\n")

	// Authz init block.
	authzBlock := ""
	if anyNeedsAuth {
		ownershipsCode := buildOwnershipsLiteral(policies)
		authzBlock = fmt.Sprintf(`
	os.Setenv("JWT_SECRET", *jwtSecret)

	if err := authz.Init(conn, %s); err != nil {
		log.Fatalf("authz init failed: %%v", err)
	}
`, ownershipsCode)
	}

	// Queue code blocks.
	queueImport := ""
	queueInitBlock := ""
	queueSubscribeBlock := ""
	if queueBackend != "" {
		subscribers := collectSubscribers(serviceFuncs)
		if len(subscribers) > 0 || hasPublishSequence(serviceFuncs) {
			queueImport = "\n\t\"context\"\n\t\"encoding/json\"\n\t\"github.com/geul-org/fullend/pkg/queue\"\n\t\"fmt\""
			queueInitBlock = fmt.Sprintf(`
	if err := queue.Init(context.Background(), %q, conn); err != nil {
		log.Fatalf("queue init failed: %%v", err)
	}
	defer queue.Close()
`, queueBackend)

			var subLines []string
			for _, fn := range subscribers {
				if fn.Param == nil {
					continue
				}
				// Determine the service package for this subscriber.
				svcPkg := "service"
				if fn.Domain != "" {
					svcPkg = fn.Domain + "svc"
				}
				subLines = append(subLines, fmt.Sprintf(`
	queue.Subscribe(%q, func(ctx context.Context, msg []byte) error {
		var message %s.%s
		if err := json.Unmarshal(msg, &message); err != nil {
			return fmt.Errorf("unmarshal: %%w", err)
		}
		return server.%s.%s(ctx, message)
	})`, fn.Subscribe.Topic, svcPkg, fn.Param.TypeName, ucFirst(fn.Domain), fn.Name))
			}
			if len(subLines) > 0 {
				queueSubscribeBlock = strings.Join(subLines, "\n") + "\n\n\tgo queue.Start(context.Background())\n"
			}
		}
	}

	// JWT flag line.
	jwtFlagLine := ""
	if anyNeedsAuth {
		jwtFlagLine = `
	jwtSecretDefault := os.Getenv("JWT_SECRET")
	if jwtSecretDefault == "" {
		jwtSecretDefault = "secret"
	}
	jwtSecret := flag.String("jwt-secret", jwtSecretDefault, "JWT signing secret")`
	}

	osImport := ""
	if anyNeedsAuth {
		osImport = "\n\t\"os\""
	}

	src := fmt.Sprintf(`package main

import (
	"database/sql"
	"flag"
	"log"%s

	_ "github.com/lib/pq"
%s%s
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	dsn := flag.String("dsn", "postgres://localhost:5432/app?sslmode=disable", "database connection string")
	dbDriver := flag.String("db", "postgres", "database driver (postgres, mysql)")%s
	flag.Parse()

	conn, err := sql.Open(*dbDriver, *dsn)
	if err != nil {
		log.Fatalf("database connection failed: %%v", err)
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		log.Fatalf("database ping failed: %%v", err)
	}
%s%s
	server := &service.Server{
%s
	}
%s
	r := service.SetupRouter(server)
	log.Printf("server listening on %%s", *addr)
	log.Fatal(r.Run(*addr))
}
`, osImport, importBlock, queueImport, jwtFlagLine, authzBlock, queueInitBlock, initBlock, queueSubscribeBlock)

	path := filepath.Join(artifactsDir, "backend", "cmd", "main.go")
	return os.WriteFile(path, []byte(src), 0644)
}
