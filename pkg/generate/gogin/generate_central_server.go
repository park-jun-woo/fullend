//ff:func feature=gen-gogin type=generator control=iteration dimension=2 topic=http-handler
//ff:what creates service/server.go that composes domain handlers with gin router

package gogin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

// generateCentralServer creates service/server.go that composes domain handlers with gin router.
func generateCentralServer(serviceDir string, domains []string, serviceFuncs []ssacparser.ServiceFunc, modulePath string, doc *openapi3.T) error {
	// Build operationId → domain map.
	opDomains := make(map[string]string)
	for _, sf := range serviceFuncs {
		if sf.Feature != "" {
			opDomains[sf.Name] = sf.Feature
		}
	}

	// Collect flat (Domain="") resources.
	flatModels := collectModelsForDomain(serviceFuncs, "")
	flatFuncs := collectFuncsForDomain(serviceFuncs, "")
	hasFlatFuncs := len(flatModels) > 0 || len(flatFuncs) > 0

	var b strings.Builder
	b.WriteString("package service\n\n")

	// Server struct.
	b.WriteString("// Server composes domain handlers.\n")
	b.WriteString("type Server struct {\n")

	// Domain handler fields.
	for _, d := range domains {
		fieldName := ucFirst(d)
		b.WriteString(fmt.Sprintf("\t%s *%ssvc.Handler\n", fieldName, d))
	}

	// Flat model fields.
	for _, m := range flatModels {
		fieldName := ucFirst(lcFirst(m) + "Model")
		b.WriteString(fmt.Sprintf("\t%s model.%sModel\n", fieldName, m))
	}
	for _, f := range flatFuncs {
		fieldName := ucFirst(f)
		b.WriteString(fmt.Sprintf("\t%s func(args ...interface{}) (interface{}, error)\n", fieldName))
	}

	// Detect security schemes from OpenAPI.
	hasBearer := hasBearerScheme(doc)

	if hasBearer {
		b.WriteString("\tJWTSecret string\n")
	}

	b.WriteString("}\n\n")

	// SetupRouter creates a gin.Engine with routes.
	b.WriteString("// SetupRouter creates a gin.Engine that routes requests to the Server.\n")
	b.WriteString("func SetupRouter(s *Server) *gin.Engine {\n")
	b.WriteString("\tr := gin.Default()\n\n")

	if hasBearer {
		b.WriteString("\t// Auth group — JWT middleware extracts currentUser into context.\n")
		b.WriteString("\tauth := r.Group(\"/\")\n")
		b.WriteString("\tauth.Use(middleware.BearerAuth(s.JWTSecret))\n\n")
	}

	writeCentralRoutes(&b, doc, opDomains, hasBearer)

	b.WriteString("\n\treturn r\n")
	b.WriteString("}\n")

	// Build imports.
	var imports []string
	imports = append(imports, "\"github.com/gin-gonic/gin\"")
	if hasBearer {
		imports = append(imports, fmt.Sprintf("\"%s/internal/middleware\"", modulePath))
	}
	if len(flatModels) > 0 || hasFlatFuncs {
		imports = append(imports, fmt.Sprintf("\"%s/internal/model\"", modulePath))
	}
	for _, d := range domains {
		imports = append(imports, fmt.Sprintf("%ssvc \"%s/internal/service/%s\"", d, modulePath, d))
	}

	var header strings.Builder
	header.WriteString("package service\n\n")
	header.WriteString("import (\n")
	for _, imp := range imports {
		header.WriteString("\t" + imp + "\n")
	}
	header.WriteString(")\n\n")

	// Replace package+import block.
	body := b.String()
	if idx := strings.Index(body, "// Server composes"); idx > 0 {
		body = body[idx:]
	}
	final := header.String() + body

	path := filepath.Join(serviceDir, "server.go")
	return os.WriteFile(path, []byte(final), 0644)
}
