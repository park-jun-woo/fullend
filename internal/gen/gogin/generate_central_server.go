//ff:func feature=gen-gogin type=generator control=iteration
//ff:what creates service/server.go that composes domain handlers with gin router

package gogin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"

	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// generateCentralServer creates service/server.go that composes domain handlers with gin router.
func generateCentralServer(serviceDir string, domains []string, serviceFuncs []ssacparser.ServiceFunc, modulePath string, doc *openapi3.T) error {
	// Build operationId → domain map.
	opDomains := make(map[string]string)
	for _, sf := range serviceFuncs {
		if sf.Domain != "" {
			opDomains[sf.Name] = sf.Domain
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

	if doc != nil {
		for pathStr, pathItem := range doc.Paths.Map() {
			for method, op := range pathItem.Operations() {
				if op.OperationID == "" {
					continue
				}
				ginPath := convertPathParamsGin(pathStr)
				handlerName := op.OperationID

				// Determine target: s.Domain.Method or s.Method.
				domain := opDomains[handlerName]
				var target string
				if domain != "" {
					target = fmt.Sprintf("s.%s.%s", ucFirst(domain), handlerName)
				} else {
					target = fmt.Sprintf("s.%s", handlerName)
				}

				// Determine route group from OpenAPI security field.
				needsAuth := opHasSecurity(op)
				ginMethod := strings.ToUpper(method)
				var routerVar string
				if needsAuth && hasBearer {
					routerVar = "auth"
				} else {
					routerVar = "r"
				}

				switch ginMethod {
				case "GET":
					b.WriteString(fmt.Sprintf("\t%s.GET(%q, %s)\n", routerVar, ginPath, target))
				case "POST":
					b.WriteString(fmt.Sprintf("\t%s.POST(%q, %s)\n", routerVar, ginPath, target))
				case "PUT":
					b.WriteString(fmt.Sprintf("\t%s.PUT(%q, %s)\n", routerVar, ginPath, target))
				case "DELETE":
					b.WriteString(fmt.Sprintf("\t%s.DELETE(%q, %s)\n", routerVar, ginPath, target))
				case "PATCH":
					b.WriteString(fmt.Sprintf("\t%s.PATCH(%q, %s)\n", routerVar, ginPath, target))
				default:
					b.WriteString(fmt.Sprintf("\t%s.Handle(%q, %q, %s)\n", routerVar, ginMethod, ginPath, target))
				}
			}
		}
	}

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
