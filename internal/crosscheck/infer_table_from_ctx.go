//ff:func feature=crosscheck type=util control=sequence topic=openapi-ddl
//ff:what operationId에서 주 테이블명 추론
package crosscheck

import (
	"github.com/getkin/kin-openapi/openapi3"

	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
)

// inferTableFromCtx guesses the primary table name from the operation's path.
func inferTableFromCtx(op *openapi3.Operation, st *ssacvalidator.SymbolTable) string {
	if op.OperationID == "" {
		return "???"
	}
	name := stripCRUDPrefix(op.OperationID)
	if name == "" {
		return "???"
	}
	table := modelToTable(name)
	if _, ok := st.DDLTables[table]; ok {
		return table
	}
	return "???"
}
