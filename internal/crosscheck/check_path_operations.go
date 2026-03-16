//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=openapi-ddl
//ff:what 단일 경로의 Operation별 x-sort/x-filter/x-include/cursor 검증
package crosscheck

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"

	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
)

// checkPathOperations validates operations for a single path.
func checkPathOperations(path string, pi *openapi3.PathItem, st *ssacvalidator.SymbolTable, funcPrimaryTable map[string]string) []CrossError {
	var errs []CrossError
	for method, op := range pi.Operations() {
		if op == nil {
			continue
		}
		ctx := fmt.Sprintf("%s %s (%s)", method, path, op.OperationID)
		primaryTable := funcPrimaryTable[op.OperationID]

		errs = append(errs, checkXSort(op, st, ctx)...)
		errs = append(errs, checkXFilter(op, st, ctx)...)
		errs = append(errs, checkXInclude(op, st, ctx, primaryTable)...)
		errs = append(errs, checkCursorSort(op, st, ctx)...)
	}
	return errs
}
