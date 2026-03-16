//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=openapi-ddl
//ff:what OpenAPI x-sort, x-filter, x-include를 DDL 테이블과 검증하고 유령 property 검사
package crosscheck

import (
	"github.com/getkin/kin-openapi/openapi3"

	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
)

// CheckOpenAPIDDL validates x-sort, x-filter, x-include against DDL tables.
func CheckOpenAPIDDL(doc *openapi3.T, st *ssacvalidator.SymbolTable, funcs []ssacparser.ServiceFunc, sensitiveCols map[string]map[string]bool) []CrossError {
	var errs []CrossError

	if doc.Paths == nil {
		return errs
	}

	funcPrimaryTable := buildFuncPrimaryTable(funcs)

	for path, pi := range doc.Paths.Map() {
		errs = append(errs, checkPathOperations(path, pi, st, funcPrimaryTable)...)
	}

	errs = append(errs, checkGhostProperties(doc, st)...)
	errs = append(errs, checkMissingProperties(doc, st, sensitiveCols)...)

	return errs
}
