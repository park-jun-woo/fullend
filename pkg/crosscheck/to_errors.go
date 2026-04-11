//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what toErrors — toulmin EvalResult를 CrossError 목록으로 변환
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func toErrors(results []toulmin.EvalResult, context string) []CrossError {
	var errs []CrossError
	for _, r := range results {
		if r.Verdict <= 0 {
			continue
		}
		switch ev := r.Evidence.(type) {
		case *rule.Evidence:
			errs = append(errs, CrossError{
				Rule:    ev.Rule,
				Context: context,
				Message: ev.Message,
				Level:   ev.Level,
			})
		case *rule.SchemaEvidence:
			errs = append(errs, CrossError{
				Rule:    ev.Rule,
				Context: context,
				Message: ev.Message,
				Level:   ev.Level,
			})
		}
	}
	return errs
}
