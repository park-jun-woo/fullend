//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what 단일 스키마의 DDL 컬럼이 OpenAPI 속성에 존재하는지 검증
package crosscheck

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"

	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
)

// checkSchemaMissingProps checks a single schema for DDL columns missing in OpenAPI.
func checkSchemaMissingProps(schemaName string, schema *openapi3.Schema, st *ssacvalidator.SymbolTable, sensitiveCols map[string]map[string]bool, xIncludeFields map[string]bool) []CrossError {
	tableName := modelToTable(schemaName)
	table, ok := st.DDLTables[tableName]
	if !ok {
		return nil
	}

	var errs []CrossError
	for colName, colType := range table.Columns {
		if isSkippedColumn(tableName, colName, sensitiveCols, xIncludeFields) {
			continue
		}
		if _, exists := schema.Properties[colName]; !exists {
			errs = append(errs, CrossError{
				Rule:       "DDL ↔ OpenAPI",
				Context:    fmt.Sprintf("table %s.%s", tableName, colName),
				Message:    fmt.Sprintf("DDL column %q (%s) — OpenAPI %s schema에 해당 property 없음", colName, colType, schemaName),
				Level:      "WARNING",
				Suggestion: fmt.Sprintf("OpenAPI %s schema에 %s property를 추가하거나, DDL에서 제거하세요", schemaName, colName),
			})
		}
	}
	return errs
}
