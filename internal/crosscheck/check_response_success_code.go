//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=ssac-openapi
//ff:what SSaC @response가 있는 함수에 OpenAPI 2xx 응답 코드가 명시되어 있는지 검증
package crosscheck

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"

	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// checkResponseSuccessCode validates that SSaC functions with @response have an explicit 2xx.
func checkResponseSuccessCode(funcs []ssacparser.ServiceFunc, doc *openapi3.T) []CrossError {
	var errs []CrossError

	opMap := buildOperationMap(doc)

	for _, fn := range funcs {
		if !hasResponseSequence(fn) {
			continue
		}

		op := opMap[fn.Name]
		if op == nil || op.Responses == nil {
			continue
		}

		if !hasExplicit2xx(op) {
			errs = append(errs, CrossError{
				Rule:       "SSaC @response → OpenAPI 2xx",
				Context:    fmt.Sprintf("%s:%s", fn.FileName, fn.Name),
				Message:    fmt.Sprintf("SSaC @response가 있는 %s에 OpenAPI 2xx 성공 응답 코드가 없습니다 (default만으로는 불충분)", fn.Name),
				Suggestion: fmt.Sprintf("OpenAPI %s responses에 200, 201, 204 등 명시적 성공 코드를 추가하세요", fn.Name),
			})
		}
	}

	return errs
}
