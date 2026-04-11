//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what validateSubscribe — @subscribe 제약 검증 (S-38~S-41, S-45)
package ssac

import (
	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

func validateSubscribe(fn parsessac.ServiceFunc) []validate.ValidationError {
	var errs []validate.ValidationError
	if fn.Param == nil || fn.Param.VarName != "message" {
		errs = append(errs, validate.ValidationError{
			Rule: "S-38", File: fn.FileName, Func: fn.Name, SeqIdx: -1, Level: "ERROR",
			Message: "@subscribe parameter variable must be named 'message'",
		})
	}
	for i, seq := range fn.Sequences {
		if seq.Type == "response" {
			errs = append(errs, validate.ValidationError{
				Rule: "S-45", File: fn.FileName, Func: fn.Name, SeqIdx: i, Level: "ERROR",
				Message: "@subscribe cannot use @response",
			})
		}
	}
	return errs
}
