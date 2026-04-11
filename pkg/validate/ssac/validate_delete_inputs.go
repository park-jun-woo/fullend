//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what validateDeleteInputs — @delete Inputs 없음 WARNING (S-11)
package ssac

import (
	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

func validateDeleteInputs(fn parsessac.ServiceFunc) []validate.ValidationError {
	var errs []validate.ValidationError
	for i, seq := range fn.Sequences {
		if seq.Type == "delete" && len(seq.Args) == 0 && !seq.SuppressWarn {
			errs = append(errs, validate.ValidationError{
				Rule: "S-11", File: fn.FileName, Func: fn.Name, SeqIdx: i, Level: "WARNING",
				Message: "@delete has no inputs — all rows may be affected",
			})
		}
	}
	return errs
}
