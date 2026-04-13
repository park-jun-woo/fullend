//ff:func feature=gen-gogin type=generator control=sequence topic=output
//ff:what domain 모드 cmd/main.go를 생성한다

package gogin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/fullend/internal/policy"
	"github.com/park-jun-woo/fullend/pkg/parser/manifest"
	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

// MainGenInput bundles generateMain's inputs to reduce parameter count.
type MainGenInput struct {
	ArtifactsDir   string
	ServiceFuncs   []ssacparser.ServiceFunc
	ModulePath     string
	QueueBackend   string
	Policies       []*policy.Policy
	SessionBackend string
	CacheBackend   string
	FileConfig     *manifest.FileBackend
}

// generateMain creates cmd/main.go with feature handler initialization.
func generateMain(in MainGenInput) error {
	artifactsDir := in.ArtifactsDir
	serviceFuncs := in.ServiceFuncs
	modulePath := in.ModulePath
	queueBackend := in.QueueBackend
	policies := in.Policies
	sessionBackend := in.SessionBackend
	cacheBackend := in.CacheBackend
	fileConfig := in.FileConfig

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
	anyNeedsAuth := anyDomainNeedsAuth(serviceFuncs, domains)

	initBlock := buildDomainInitBlock(serviceFuncs, domains, anyNeedsAuth)
	importBlock := buildDomainImportsBlock(domains, modulePath, anyNeedsAuth)

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

	queueImport, queueInitBlock, queueSubscribeBlock := buildQueueBlocks(serviceFuncs, queueBackend)

	builtinImport, builtinInitBlock := buildBuiltinInitBlocks(sessionBackend, cacheBackend, fileConfig)
	if strings.Contains(queueImport, `"context"`) {
		builtinImport = strings.Replace(builtinImport, "\n\t\"context\"", "", 1)
	}

	jwtFlagLine := ""
	osImport := ""
	if anyNeedsAuth {
		jwtFlagLine = `
	jwtSecretDefault := os.Getenv("JWT_SECRET")
	if jwtSecretDefault == "" {
		jwtSecretDefault = "secret"
	}
	jwtSecret := flag.String("jwt-secret", jwtSecretDefault, "JWT signing secret")`
		osImport = "\n\t\"os\""
	}

	src := mainWithDomainsTemplate(osImport, importBlock, queueImport, builtinImport, jwtFlagLine, authzBlock, queueInitBlock, builtinInitBlock, initBlock, queueSubscribeBlock)

	path := filepath.Join(artifactsDir, "backend", "cmd", "main.go")
	return os.WriteFile(path, []byte(src), 0644)
}
