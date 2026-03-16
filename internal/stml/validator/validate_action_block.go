//ff:func feature=stml-validate type=rule control=iteration dimension=1
//ff:what data-action 블록의 operationId·메서드·필드 검증
package validator

import (
	"fmt"
	"strings"

	"github.com/geul-org/fullend/internal/stml/parser"
)

func validateActionBlock(a parser.ActionBlock, file string, st *SymbolTable, frontendDir string) []ValidationError {
	var errs []ValidationError
	attr := fmt.Sprintf("data-action=%q", a.OperationID)

	api, ok := st.Operations[a.OperationID]
	if !ok {
		return append(errs, errOpNotFound(file, attr, a.OperationID))
	}

	if api.Method == "GET" {
		errs = append(errs, errWrongMethod(file, attr, a.OperationID, api.Method, "POST/PUT/DELETE"))
	}

	errs = append(errs, validateParams(a.Params, a.OperationID, file, api)...)

	for _, f := range a.Fields {
		if strings.HasPrefix(f.Tag, "data-component:") {
			comp := strings.TrimPrefix(f.Tag, "data-component:")
			errs = append(errs, validateComponent(comp, file, frontendDir)...)
		}
		if _, ok := api.RequestFields[f.Name]; !ok {
			errs = append(errs, errFieldNotFound(file, a.OperationID, f.Name))
		}
	}

	return errs
}
