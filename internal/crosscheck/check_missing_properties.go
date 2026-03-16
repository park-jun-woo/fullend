//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=openapi-ddl
//ff:what DDL 컬럼이 OpenAPI 스키마에 property로 존재하는지 검증
package crosscheck

import (
	"github.com/getkin/kin-openapi/openapi3"

	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
)

// checkMissingProperties checks DDL columns -> OpenAPI schema properties.
func checkMissingProperties(doc *openapi3.T, st *ssacvalidator.SymbolTable, sensitiveCols map[string]map[string]bool) []CrossError {
	var errs []CrossError
	if doc.Components == nil || doc.Components.Schemas == nil {
		return errs
	}

	xIncludeFields := collectXIncludeLocalFields(doc)

	for schemaName, schemaRef := range doc.Components.Schemas {
		if schemaRef == nil || schemaRef.Value == nil {
			continue
		}
		errs = append(errs, checkSchemaMissingProps(schemaName, schemaRef.Value, st, sensitiveCols, xIncludeFields)...)
	}
	return errs
}
