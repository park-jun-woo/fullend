//ff:func feature=ssac-validate type=rule control=iteration dimension=2
//ff:what HTTP 함수에서 message 사용을 검증한다
package validator

import (
	"strings"

	"github.com/geul-org/fullend/internal/ssac/parser"
)

// validateNoMessageInHTTP는 HTTP 함수에서 message 사용을 검증한다.
func validateNoMessageInHTTP(sf parser.ServiceFunc) []ValidationError {
	var errs []ValidationError
	for i, seq := range sf.Sequences {
		for _, val := range seq.Inputs {
			if strings.HasPrefix(val, "message.") {
				ctx := errCtx{sf.FileName, sf.Name, i}
				errs = append(errs, ctx.err("@sequence", "HTTP 함수에서 message를 사용할 수 없습니다 — @subscribe 함수에서만 사용 가능합니다"))
				break
			}
		}
	}
	return errs
}
