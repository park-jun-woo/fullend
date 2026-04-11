//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what checkStaleFields — @response에서 mutated 후 refresh 안 된 변수 검출
package ssac

import (
	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

func checkStaleFields(fn parsessac.ServiceFunc, seq parsessac.Sequence, mutated, refreshed map[string]bool) []validate.ValidationError {
	if seq.SuppressWarn {
		return nil
	}
	var errs []validate.ValidationError
	for _, varRef := range seq.Fields {
		varType := findVarType(fn.Sequences, varRef)
		if varType != "" && mutated[varType] && !refreshed[varType] {
			errs = append(errs, validate.ValidationError{
				Rule: "S-36", File: fn.FileName, Func: fn.Name, SeqIdx: -1, Level: "WARNING",
				Message: "@response uses " + varRef + " which was mutated but not re-queried",
			})
		}
	}
	return errs
}
