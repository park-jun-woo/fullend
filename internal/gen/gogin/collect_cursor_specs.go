//ff:func feature=gen-gogin type=util control=iteration dimension=2
//ff:what extracts cursor column Go field name per operationId from OpenAPI x-pagination

package gogin

import (
	"github.com/ettle/strcase"
	"github.com/getkin/kin-openapi/openapi3"
)

// collectCursorSpecs extracts cursor column Go field name per operationId.
// Returns map[operationId]string ("ID" default, or PascalCase of x-sort.default).
func collectCursorSpecs(doc *openapi3.T) map[string]string {
	result := make(map[string]string)
	if doc == nil || doc.Paths == nil {
		return result
	}
	for _, pi := range doc.Paths.Map() {
		for _, op := range pi.Operations() {
			if op == nil || op.OperationID == "" {
				continue
			}
			pag := getExtMap(op, "x-pagination")
			if pag == nil || getStr(pag, "style", "") != "cursor" {
				continue
			}
			cursorField := "ID"
			def := ""
			if sortExt := getExtMap(op, "x-sort"); sortExt != nil {
				def = getStr(sortExt, "default", "")
			}
			if def != "" {
				cursorField = strcase.ToGoPascal(def)
			}
			result[op.OperationID] = cursorField
		}
	}
	return result
}
