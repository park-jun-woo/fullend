//ff:func feature=ssac-validate type=rule control=sequence topic=subscribe
//ff:what subscribe/HTTP 트리거 관련 규칙 검증
package validator

import (
	"github.com/geul-org/fullend/internal/ssac/parser"
)

// validateSubscribeRules는 subscribe/HTTP 트리거와 관련된 규칙을 검증한다.
func validateSubscribeRules(sf parser.ServiceFunc) []ValidationError {
	if sf.Subscribe != nil {
		return validateSubscribeConstraints(sf)
	}
	return validateNoMessageInHTTP(sf)
}
