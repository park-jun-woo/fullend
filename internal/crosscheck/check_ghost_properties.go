//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what OpenAPI 스키마 속성이 DDL 컬럼에 대응하는지 검증 (유령 property 탐지)
package crosscheck

import (
	"github.com/getkin/kin-openapi/openapi3"

	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
)

// checkGhostProperties detects OpenAPI schema properties that have no corresponding DDL column.
func checkGhostProperties(doc *openapi3.T, st *ssacvalidator.SymbolTable) []CrossError {
	var errs []CrossError

	if doc.Components == nil || doc.Components.Schemas == nil {
		return errs
	}

	xIncludeFields := collectXIncludeLocalFields(doc)

	for schemaName, schemaRef := range doc.Components.Schemas {
		if schemaRef == nil || schemaRef.Value == nil {
			continue
		}
		errs = append(errs, checkSchemaGhostProps(schemaName, schemaRef.Value, st, xIncludeFields)...)
	}

	return errs
}
