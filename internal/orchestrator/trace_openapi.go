//ff:func feature=orchestrator type=util
//ff:what traceOpenAPI finds the OpenAPI path/method for an operationId.

package orchestrator

import (
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func traceOpenAPI(doc *openapi3.T, opID string, specsDir string) *ChainLink {
	if doc.Paths == nil {
		return nil
	}
	for path, pi := range doc.Paths.Map() {
		for method, op := range pi.Operations() {
			if op.OperationID == opID {
				line := grepLine(filepath.Join(specsDir, "api", "openapi.yaml"), "operationId: "+opID)
				return &ChainLink{
					Kind:    "OpenAPI",
					File:    "api/openapi.yaml",
					Line:    line,
					Summary: strings.ToUpper(method) + " " + path,
				}
			}
		}
	}
	return nil
}
