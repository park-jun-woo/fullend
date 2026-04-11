//ff:func feature=rule type=util control=iteration dimension=1
//ff:what toValidationErrors — toulmin EvalResult를 ValidationError로 변환 (ssac 내부용)
package ssac

import (
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func toValidationErrors(results []toulmin.EvalResult, file, funcName string, seqIdx int) []validate.ValidationError {
	var errs []validate.ValidationError
	for _, r := range results {
		if r.Verdict <= 0 {
			continue
		}
		if ev, ok := r.Evidence.(*rule.Evidence); ok {
			errs = append(errs, validate.ValidationError{
				Rule: ev.Rule, File: file, Func: funcName,
				SeqIdx: seqIdx, Level: ev.Level, Message: ev.Message,
			})
		}
	}
	return errs
}
