//ff:func feature=rule type=util control=iteration dimension=1
//ff:what toSTMLErrors — toulmin EvalResult를 STML ValidationError로 변환
package stml

import (
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func toSTMLErrors(results []toulmin.EvalResult, file, context string) []validate.ValidationError {
	var errs []validate.ValidationError
	for _, r := range results {
		if r.Verdict <= 0 {
			continue
		}
		if ev, ok := r.Evidence.(*rule.Evidence); ok {
			errs = append(errs, validate.ValidationError{
				Rule: ev.Rule, File: file, Func: context,
				SeqIdx: -1, Level: ev.Level, Message: ev.Message,
			})
		}
	}
	return errs
}
