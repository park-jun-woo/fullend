//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=config-check
//ff:what 보안 요구사항의 각 이름이 미들웨어에 존재하는지 검증
package crosscheck

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func checkSecurityReqNames(req openapi3.SecurityRequirement, mwSet map[string]bool, method, pathStr string) []CrossError {
	var errs []CrossError
	for name := range req {
		if !mwSet[name] {
			errs = append(errs, CrossError{
				Rule:       "Config ↔ OpenAPI",
				Context:    fmt.Sprintf("%s %s", strings.ToUpper(method), pathStr),
				Message:    fmt.Sprintf("endpoint references security %q not in fullend.yaml middleware", name),
				Suggestion: fmt.Sprintf("fullend.yaml backend.middleware에 %q 추가", name),
			})
		}
	}
	return errs
}
