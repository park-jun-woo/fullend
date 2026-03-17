//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=ssac-openapi
//ff:what SSaC 함수명과 OpenAPI operationId 교차 검증 및 @response 필드 매칭
package crosscheck

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/park-jun-woo/fullend/internal/funcspec"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
)

// CheckSSaCOpenAPI validates SSaC function names match OpenAPI operationIds and vice versa,
// and SSaC @response fields match OpenAPI response schema properties.
func CheckSSaCOpenAPI(funcs []ssacparser.ServiceFunc, st *ssacvalidator.SymbolTable, doc *openapi3.T, funcSpecs []funcspec.FuncSpec) []CrossError {
	var errs []CrossError

	funcNames := make(map[string]string)
	for _, fn := range funcs {
		if fn.Subscribe != nil {
			continue // @subscribe는 HTTP endpoint가 아니므로 operationId 불필요
		}
		funcNames[fn.Name] = fn.FileName
	}

	for name, fileName := range funcNames {
		if _, ok := st.Operations[name]; !ok {
			errs = append(errs, CrossError{
				Rule:       "SSaC → OpenAPI",
				Context:    fmt.Sprintf("%s:%s", fileName, name),
				Message:    fmt.Sprintf("SSaC function %q has no matching OpenAPI operationId", name),
				Suggestion: fmt.Sprintf("OpenAPI에 추가: operationId: %s", name),
			})
		}
	}

	for opID := range st.Operations {
		if _, ok := funcNames[opID]; !ok {
			errs = append(errs, CrossError{
				Rule:       "OpenAPI → SSaC",
				Context:    opID,
				Message:    fmt.Sprintf("OpenAPI operationId %q has no matching SSaC function", opID),
				Level:      "WARNING",
				Suggestion: fmt.Sprintf("SSaC에 추가: func %s(w http.ResponseWriter, r *http.Request) {}", opID),
			})
		}
	}

	errs = append(errs, checkResponseFields(funcs, st, doc, funcSpecs)...)

	if doc != nil {
		errs = append(errs, checkErrStatus(funcs, doc)...)
	}

	if doc != nil {
		errs = append(errs, checkResponseSuccessCode(funcs, doc)...)
	}

	return errs
}
