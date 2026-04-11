//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what validateErrStatus — IANA 미등록 HTTP status 검증 (S-58)
package ssac

import (
	"fmt"

	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

func validateErrStatus(fn parsessac.ServiceFunc) []validate.ValidationError {
	var errs []validate.ValidationError
	for i, seq := range fn.Sequences {
		if seq.ErrStatus > 0 && (seq.ErrStatus < 100 || seq.ErrStatus > 599) {
			errs = append(errs, validate.ValidationError{
				Rule: "S-58", File: fn.FileName, Func: fn.Name, SeqIdx: i, Level: "ERROR",
				Message: fmt.Sprintf("ErrStatus %d is not a valid HTTP status code (100-599)", seq.ErrStatus),
			})
		}
	}
	return errs
}
