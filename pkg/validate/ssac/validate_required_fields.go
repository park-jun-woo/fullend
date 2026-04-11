//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what validateRequiredFields — 시퀀스 타입별 필수/금지 필드 검증 (S-1~S-24)
package ssac

import (
	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

func validateRequiredFields(fn parsessac.ServiceFunc, ground *rule.Ground) []validate.ValidationError {
	var errs []validate.ValidationError
	for i, seq := range fn.Sequences {
		errs = append(errs, validateSeqFields(fn.FileName, fn.Name, i, seq, ground)...)
	}
	return errs
}
