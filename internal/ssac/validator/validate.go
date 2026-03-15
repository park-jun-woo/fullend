//ff:func feature=ssac-validate type=command
//ff:what SSaC 검증 메인 엔트리포인트 — []ServiceFunc의 내부 정합성을 검증한다
package validator

import (
	"github.com/geul-org/fullend/internal/ssac/parser"
)

// Validate는 []ServiceFunc의 내부 정합성을 검증한다.
func Validate(funcs []parser.ServiceFunc) []ValidationError {
	var errs []ValidationError
	for _, sf := range funcs {
		errs = append(errs, validateFunc(sf)...)
	}
	return errs
}
