//ff:func feature=ssac-validate type=rule control=iteration dimension=1 topic=args-inputs
//ff:what ErrStatus가 IANA 등록 HTTP 상태 코드인지 검증
package validator

import (
	"fmt"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

// validateErrStatus validates that ErrStatus values are IANA-registered HTTP status codes.
func validateErrStatus(sf parser.ServiceFunc) []ValidationError {
	var errs []ValidationError
	for i, seq := range sf.Sequences {
		if seq.ErrStatus == 0 {
			continue
		}
		if !validHTTPStatus[seq.ErrStatus] {
			ctx := errCtx{sf.FileName, sf.Name, i}
			errs = append(errs, ctx.warn("@"+seq.Type, fmt.Sprintf("HTTP status %d는 IANA 등록 코드가 아닙니다", seq.ErrStatus)))
		}
	}
	return errs
}
