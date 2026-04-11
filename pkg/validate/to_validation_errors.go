//ff:func feature=rule type=util control=iteration dimension=1
//ff:what toValidationErrors — toulmin EvalResult를 ValidationError 목록으로 변환
package validate

import (
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func toValidationErrors(results []toulmin.EvalResult, file, funcName string, seqIdx int) []ValidationError {
	var errs []ValidationError
	for _, r := range results {
		if r.Verdict <= 0 {
			continue
		}
		if ev, ok := r.Evidence.(*rule.Evidence); ok {
			errs = append(errs, ValidationError{
				Rule: ev.Rule, File: file, Func: funcName,
				SeqIdx: seqIdx, Level: ev.Level, Message: ev.Message,
			})
		}
	}
	return errs
}
