//ff:func feature=gen-gogin type=generator control=iteration dimension=2
//ff:what creates service/server.go with Server struct definition and Handler function

package gogin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// generateServerStruct creates service/server.go with Server struct definition and Handler function.
func generateServerStruct(intDir string, models, funcs []string, modulePath string, doc *openapi3.T) error {
	serviceDir := filepath.Join(intDir, "service")
	if err := os.MkdirAll(serviceDir, 0755); err != nil {
		return err
	}

	var b strings.Builder

	b.WriteString("package service\n\n")

	// Import model package if there are models.
	if len(models) > 0 {
		b.WriteString(fmt.Sprintf("import \"%s/internal/model\"\n\n", modulePath))
	}

	// Build fields.
	var fields []string
	for _, m := range models {
		fieldName := ucFirst(lcFirst(m) + "Model")
		fields = append(fields, fmt.Sprintf("\t%s model.%sModel", fieldName, m))
	}
	for _, f := range funcs {
		fieldName := ucFirst(f)
		fields = append(fields, fmt.Sprintf("\t%s func(args ...interface{}) (interface{}, error)", fieldName))
	}
	b.WriteString("// Server implements api.ServerInterface.\n")
	b.WriteString("type Server struct {\n")
	for _, f := range fields {
		b.WriteString(f + "\n")
	}
	b.WriteString("}\n\n")

	// Handler function.
	b.WriteString("// Handler creates an http.Handler that routes requests to the Server.\n")
	b.WriteString("func Handler(s *Server) http.Handler {\n")
	b.WriteString("\tmux := http.NewServeMux()\n")

	for pathStr, pathItem := range doc.Paths.Map() {
		for method, op := range pathItem.Operations() {
			if op.OperationID == "" {
				continue
			}
			muxPath := convertPathParams(pathStr)
			pattern := fmt.Sprintf("%s %s", method, muxPath)
			handlerName := op.OperationID
			pathParams := collectPathParams(pathItem, op)

			if len(pathParams) == 0 {
				b.WriteString(fmt.Sprintf("\tmux.HandleFunc(\"%s\", s.%s)\n", pattern, handlerName))
			} else {
				writeRouteHandler(&b, pattern, handlerName, pathParams)
			}
		}
	}

	b.WriteString("\treturn mux\n")
	b.WriteString("}\n")

	// Add imports at the top based on content.
	content := b.String()
	imports := []string{fmt.Sprintf("\"%s/internal/model\"", modulePath)}
	if strings.Contains(content, "http.") {
		imports = append(imports, "\"net/http\"")
	}
	if strings.Contains(content, "strconv.") {
		imports = append(imports, "\"strconv\"")
	}

	var header strings.Builder
	header.WriteString("package service\n\n")
	header.WriteString("import (\n")
	for _, imp := range imports {
		header.WriteString("\t" + imp + "\n")
	}
	header.WriteString(")\n\n")

	// Replace the original "package service\n\nimport ..." section.
	// Remove the original package + import block and prepend the new one.
	body := content
	if idx := strings.Index(body, "// Server implements"); idx > 0 {
		body = body[idx:]
	}
	final := header.String() + body

	path := filepath.Join(serviceDir, "server.go")
	return os.WriteFile(path, []byte(final), 0644)
}
