//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what validateConfigRef — SSaC에서 config.* 사용 금지 검증 (S-31)
package ssac

import (

	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

func validateConfigRef(fn parsessac.ServiceFunc) []validate.ValidationError {
	var errs []validate.ValidationError
	for i, seq := range fn.Sequences {
		errs = append(errs, checkConfigInputs(fn.FileName, fn.Name, i, seq.Inputs)...)
	}
	return errs
}
