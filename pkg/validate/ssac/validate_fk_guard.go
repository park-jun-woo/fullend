//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what validateFKGuard — FK 참조 @get 후 @empty 가드 누락 검증 (S-37)
package ssac

import (

	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

func validateFKGuard(fn parsessac.ServiceFunc) []validate.ValidationError {
	declared := map[string]bool{}
	varTypes := map[string]string{}
	if fn.Subscribe != nil {
		declared["message"] = true
	}
	var errs []validate.ValidationError
	for i, seq := range fn.Sequences {
		errs = append(errs, checkFKGet(fn, i, seq, declared, varTypes)...)
		if seq.Result != nil {
			declared[seq.Result.Var] = true
			varTypes[seq.Result.Var] = seq.Result.Type
		}
	}
	return errs
}
