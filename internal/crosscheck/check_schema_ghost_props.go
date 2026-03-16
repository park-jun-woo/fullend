//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what 단일 스키마의 속성이 DDL 컬럼에 존재하는지 검증
package crosscheck

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"

	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
)

// checkSchemaGhostProps checks a single schema's properties against DDL columns.
func checkSchemaGhostProps(schemaName string, schema *openapi3.Schema, st *ssacvalidator.SymbolTable, xIncludeFields map[string]bool) []CrossError {
	tableName := modelToTable(schemaName)
	table, ok := st.DDLTables[tableName]
	if !ok {
		return nil
	}

	var errs []CrossError
	for propName := range schema.Properties {
		if xIncludeFields[propName] {
			continue
		}
		if _, colExists := table.Columns[propName]; !colExists {
			errs = append(errs, CrossError{
				Rule:       "OpenAPI ↔ DDL",
				Context:    fmt.Sprintf("schema %s.%s", schemaName, propName),
				Message:    fmt.Sprintf("OpenAPI property %q — DDL %s에 해당 컬럼 없음 (유령 property)", propName, tableName),
				Level:      "ERROR",
				Suggestion: fmt.Sprintf("DDL에 추가하거나 OpenAPI에서 제거: %s.%s", tableName, propName),
			})
		}
	}
	return errs
}
