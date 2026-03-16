//ff:func feature=stml-validate type=rule control=iteration dimension=1
//ff:what data-fetch 블록의 operationId·메서드·파라미터·바인딩 검증
package validator

import (
	"fmt"

	"github.com/geul-org/fullend/internal/stml/parser"
)

func validateFetchBlock(f parser.FetchBlock, file string, st *SymbolTable, cs *CustomSymbol, frontendDir string) []ValidationError {
	var errs []ValidationError
	attr := fmt.Sprintf("data-fetch=%q", f.OperationID)

	api, ok := st.Operations[f.OperationID]
	if !ok {
		return append(errs, errOpNotFound(file, attr, f.OperationID))
	}

	if api.Method != "GET" {
		errs = append(errs, errWrongMethod(file, attr, f.OperationID, api.Method, "GET"))
	}

	errs = append(errs, validateParams(f.Params, f.OperationID, file, api)...)
	errs = append(errs, validateFetchBinds(f.Binds, f.OperationID, file, api, cs)...)
	errs = append(errs, validateFetchEaches(f.Eaches, f.OperationID, file, api)...)

	for _, c := range f.Components {
		errs = append(errs, validateComponent(c.Name, file, frontendDir)...)
	}

	errs = append(errs, validateInfraParams(f, file, api)...)
	errs = append(errs, validateNestedFetches(f, file, st, cs, frontendDir)...)
	errs = append(errs, validateChildActions(f, file, st, frontendDir)...)

	return errs
}
