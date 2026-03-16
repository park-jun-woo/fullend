//ff:func feature=gen-gogin type=generator control=sequence topic=output
//ff:what domain 모드 cmd/main.go를 생성한다

package gogin

import (
	"fmt"
	"os"
	"path/filepath"

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

	src := mainWithDomainsTemplate(osImport, importBlock, queueImport, jwtFlagLine, authzBlock, queueInitBlock, initBlock, queueSubscribeBlock)

	path := filepath.Join(artifactsDir, "backend", "cmd", "main.go")
	return os.WriteFile(path, []byte(src), 0644)
}
