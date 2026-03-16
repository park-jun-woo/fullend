//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=ssac-openapi
//ff:what shorthand @response 변수를 OpenAPI 속성과 대조 검증
package crosscheck

import (
	"fmt"

	"github.com/geul-org/fullend/internal/funcspec"
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
)

// checkShorthandResponse validates shorthand @response variables against OpenAPI.
func checkShorthandResponse(fn ssacparser.ServiceFunc, funcSpecs []funcspec.FuncSpec, st *ssacvalidator.SymbolTable, opResponseProps map[string]map[string]bool) []CrossError {
	shorthandFields := resolveShorthandResponseFields(fn, funcSpecs, st)
	if shorthandFields == nil {
		return nil
	}

	opProps, hasOp := opResponseProps[fn.Name]
	if !hasOp {
		return nil
	}

	var errs []CrossError
	shorthandSet := make(map[string]bool, len(shorthandFields))
	for _, f := range shorthandFields {
		shorthandSet[f] = true
	}

	for _, jf := range shorthandFields {
		if !opProps[jf] {
			errs = append(errs, CrossError{
				Rule:       "SSaC @response → OpenAPI",
				Context:    fmt.Sprintf("%s:%s", fn.FileName, fn.Name),
				Message:    fmt.Sprintf("shorthand @response 변수의 JSON 필드 %q가 OpenAPI %s 응답 스키마에 없습니다", jf, fn.Name),
				Suggestion: fmt.Sprintf("OpenAPI %s 응답 스키마의 property명을 %q로 변경하세요", fn.Name, jf),
			})
		}
	}

	for prop := range opProps {
		if !shorthandSet[prop] {
			errs = append(errs, CrossError{
				Rule:       "OpenAPI → SSaC @response",
				Context:    fmt.Sprintf("%s:%s", fn.FileName, fn.Name),
				Message:    fmt.Sprintf("OpenAPI %s 응답 필드 %q가 shorthand @response 변수 타입에 없습니다", fn.Name, prop),
				Level:      "WARNING",
				Suggestion: fmt.Sprintf("OpenAPI에서 %q를 제거하거나 변수 타입에 해당 필드를 추가하세요", prop),
			})
		}
	}

	return errs
}
