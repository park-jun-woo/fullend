//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what validateMethodRef — Model.Method의 Method 존재 검증 (S-49)
package ssac

import (
	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

func validateMethodRef(fn parsessac.ServiceFunc, ground *rule.Ground) []validate.ValidationError {
	var errs []validate.ValidationError
	for i, seq := range fn.Sequences {
		errs = append(errs, checkMethodRefSeq(ground, fn.FileName, fn.Name, i, seq)...)
	}
	return errs
}
