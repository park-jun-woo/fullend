//ff:func feature=gen-gogin type=generator control=iteration dimension=2
//ff:what creates backend/go.mod and backend/cmd/main.go for flat mode

package gogin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/geul-org/fullend/internal/policy"
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// generateMain creates backend/go.mod (if missing) and backend/cmd/main.go.
func generateMain(artifactsDir string, models []string, modulePath string, queueBackend string, serviceFuncs []ssacparser.ServiceFunc, policies []*policy.Policy) error {
	if modulePath == "" {
		base := filepath.Base(artifactsDir)
		modulePath = base + "/backend"
	}

	// Generate backend/go.mod if it doesn't exist.
	goModPath := filepath.Join(artifactsDir, "backend", "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Join(artifactsDir, "backend"), 0755); err != nil {
			return err
		}
		goModContent := fmt.Sprintf("module %s\n\ngo 1.22\n", modulePath)
		if err := os.WriteFile(goModPath, []byte(goModContent), 0644); err != nil {
			return err
		}
	}

	if err := os.MkdirAll(filepath.Join(artifactsDir, "backend", "cmd"), 0755); err != nil {
		return err
	}

	// Build model field initialization lines.
	var initLines []string
	for _, m := range models {
		fieldName := ucFirst(lcFirst(m) + "Model")
		initLines = append(initLines, fmt.Sprintf("\t\t%s: model.New%sModel(conn),", fieldName, m))
	}
	initBlock := strings.Join(initLines, "\n")
	if initBlock == "" {
		initBlock = "\t\t// No models detected"
	}

	// Authz init block.
	authzImport := ""
	authzInitBlock := ""
	if hasAuthSequence(serviceFuncs) {
		authzImport = "\n\t\"github.com/geul-org/fullend/pkg/authz\""
		ownershipsCode := buildOwnershipsLiteral(policies)
		authzInitBlock = fmt.Sprintf(`
	if err := authz.Init(conn, %s); err != nil {
		log.Fatalf("authz init failed: %%v", err)
	}
`, ownershipsCode)
	}

	// Queue code blocks.
	queueImport := ""
	queueInitBlock := ""
	queueSubscribeBlock := ""
	subscribers := collectSubscribers(serviceFuncs)
	needsQueue := queueBackend != "" && (len(subscribers) > 0 || hasPublishSequence(serviceFuncs))
	if needsQueue {
		queueImport = "\n\t\"context\"\n\t\"encoding/json\"\n\t\"github.com/geul-org/fullend/pkg/queue\""
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
			subLines = append(subLines, fmt.Sprintf(`
	queue.Subscribe(%q, func(ctx context.Context, msg []byte) error {
		var message service.%s
		if err := json.Unmarshal(msg, &message); err != nil {
			return fmt.Errorf("unmarshal: %%w", err)
		}
		return server.%s(ctx, message)
	})`, fn.Subscribe.Topic, fn.Param.TypeName, fn.Name))
		}
		if len(subLines) > 0 {
			queueSubscribeBlock = strings.Join(subLines, "\n") + "\n\n\tgo queue.Start(context.Background())\n"
			queueImport += "\n\t\"fmt\""
		}
	}

	src := mainTemplate(modulePath, authzImport, queueImport, authzInitBlock, queueInitBlock, initBlock, queueSubscribeBlock)

	path := filepath.Join(artifactsDir, "backend", "cmd", "main.go")
	return os.WriteFile(path, []byte(src), 0644)
}
