//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what SSaC @response 필드 키가 OpenAPI 응답 스키마 속성과 일치하는지 검증
package crosscheck

import (
	"github.com/getkin/kin-openapi/openapi3"

	"github.com/geul-org/fullend/internal/funcspec"
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
)

// checkResponseFields validates that SSaC @response field keys match OpenAPI response schema properties.
func checkResponseFields(funcs []ssacparser.ServiceFunc, st *ssacvalidator.SymbolTable, doc *openapi3.T, funcSpecs []funcspec.FuncSpec) []CrossError {
	var errs []CrossError

	opResponseProps := buildOperationResponseProps(doc)

	for _, fn := range funcs {
		responseFields := extractResponseFieldKeys(fn)

		if responseFields == nil {
			errs = append(errs, checkShorthandResponse(fn, funcSpecs, st, opResponseProps)...)
			continue
		}

		opProps, hasOp := opResponseProps[fn.Name]
		if !hasOp {
			continue
		}

		errs = append(errs, checkExplicitResponseFields(fn, responseFields, opProps)...)
	}

	return errs
}
