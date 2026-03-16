//ff:func feature=gen-gogin type=util control=iteration dimension=4
//ff:what extracts per-model x-include specs from OpenAPI operations

package gogin

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"

	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// collectModelIncludes extracts per-model x-include specs from OpenAPI operations.
// Maps operationId → model via serviceFuncs, then merges all include specs per model.
// Returns map[ModelName][]string where values are raw include specs (e.g. "user", "instructor_id:user").
func collectModelIncludes(doc *openapi3.T, funcs []ssacparser.ServiceFunc) map[string][]string {
	// Map operationId → model name (from the first get sequence with @model).
	opToModel := make(map[string]string)
	for _, fn := range funcs {
		for _, seq := range fn.Sequences {
			if seq.Model != "" && seq.Type == "get" {
				parts := strings.SplitN(seq.Model, ".", 2)
				if len(parts) >= 1 {
					opToModel[fn.Name] = parts[0]
					break
				}
			}
		}
	}

	result := make(map[string][]string)
	if doc == nil || doc.Paths == nil {
		return result
	}

	for _, pathItem := range doc.Paths.Map() {
		for _, op := range pathItem.Operations() {
			if op.OperationID == "" {
				continue
			}
			incCfg := getExtMap(op, "x-include")
			if incCfg == nil {
				continue
			}
			allowed := getStrSlice(incCfg, "allowed")
			if len(allowed) == 0 {
				continue
			}
			modelName, ok := opToModel[op.OperationID]
			if !ok {
				continue
			}
			existing := result[modelName]
			for _, spec := range allowed {
				found := false
				for _, e := range existing {
					if e == spec {
						found = true
						break
					}
				}
				if !found {
					existing = append(existing, spec)
				}
			}
			result[modelName] = existing
		}
	}

	return result
}
