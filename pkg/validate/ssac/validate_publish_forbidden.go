//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what validatePublishForbidden — @publish에서 query 사용 금지 (S-32)
package ssac

import (
	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

func validatePublishForbidden(fn parsessac.ServiceFunc) []validate.ValidationError {
	var errs []validate.ValidationError
	for i, seq := range fn.Sequences {
		if seq.Type != "publish" {
			continue
		}
		errs = append(errs, checkPublishQueryForbidden(fn.FileName, fn.Name, i, seq.Inputs)...)
	}
	return errs
}
