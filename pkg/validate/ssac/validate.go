//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what Validate — SSaC ServiceFunc 전체 검증 (S-1~S-58)
package ssac

import (
	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

// Validate checks all SSaC service functions.
func Validate(funcs []parsessac.ServiceFunc, ground *rule.Ground) []validate.ValidationError {
	var errs []validate.ValidationError
	for _, fn := range funcs {
		errs = append(errs, validateFunc(fn, ground)...)
	}
	return errs
}
