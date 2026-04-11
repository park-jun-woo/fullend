//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what validateSubscribe — @subscribe 제약 검증 (S-38~S-41, S-45)
package ssac

import (
	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

func validateSubscribe(fn parsessac.ServiceFunc) []validate.ValidationError {
	var errs []validate.ValidationError
	// S-38: parameter variable must be "message"
	if fn.Param == nil || fn.Param.VarName != "message" {
		errs = append(errs, validate.ValidationError{
			Rule: "S-38", File: fn.FileName, Func: fn.Name, SeqIdx: -1, Level: "ERROR",
			Message: "@subscribe parameter variable must be named 'message'",
		})
	}
	// S-39: message type must have matching struct
	if fn.Param != nil && fn.Param.TypeName != "" {
		if !hasStructDef(fn.Structs, fn.Param.TypeName) {
			errs = append(errs, validate.ValidationError{
				Rule: "S-39", File: fn.FileName, Func: fn.Name, SeqIdx: -1, Level: "ERROR",
				Message: "@subscribe message type " + fn.Param.TypeName + " has no struct definition",
			})
		}
	}
	// S-40: @subscribe cannot use request
	// S-41: @subscribe cannot use query
	for i, seq := range fn.Sequences {
		errs = append(errs, checkSubscribeSeq(fn.FileName, fn.Name, i, seq)...)
	}
	return errs
}
