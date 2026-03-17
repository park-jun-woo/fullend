//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=ssac-openapi
//ff:what 명시적 @response 필드를 OpenAPI 속성과 대조 검증
package crosscheck

import (
	"fmt"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

// checkExplicitResponseFields validates explicit @response fields against OpenAPI.
func checkExplicitResponseFields(fn ssacparser.ServiceFunc, responseFields []string, opProps map[string]bool) []CrossError {
	var errs []CrossError

	for _, field := range responseFields {
		if !opProps[field] {
			errs = append(errs, CrossError{
				Rule:       "SSaC @response → OpenAPI",
				Context:    fmt.Sprintf("%s:%s", fn.FileName, fn.Name),
				Message:    fmt.Sprintf("SSaC @response 필드 %q가 OpenAPI %s 응답 스키마에 없습니다", field, fn.Name),
				Suggestion: fmt.Sprintf("OpenAPI %s 응답 스키마에 %q property를 추가하세요", fn.Name, field),
			})
		}
	}

	responseFieldSet := make(map[string]bool, len(responseFields))
	for _, f := range responseFields {
		responseFieldSet[f] = true
	}
	for prop := range opProps {
		if !responseFieldSet[prop] {
			errs = append(errs, CrossError{
				Rule:       "OpenAPI → SSaC @response",
				Context:    fmt.Sprintf("%s:%s", fn.FileName, fn.Name),
				Message:    fmt.Sprintf("OpenAPI %s 응답 필드 %q가 SSaC @response에 없습니다", fn.Name, prop),
				Level:      "WARNING",
				Suggestion: fmt.Sprintf("SSaC @response에 %q 필드를 추가하거나 OpenAPI에서 제거하세요", prop),
			})
		}
	}

	return errs
}
