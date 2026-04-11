//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what checkSubscribeSeq — @subscribe 시퀀스에서 request/query/@response 사용 검증 (S-40, S-41, S-45)
package ssac

import (
	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

func checkSubscribeSeq(file, funcName string, seqIdx int, seq parsessac.Sequence) []validate.ValidationError {
	var errs []validate.ValidationError
	if seq.Type == "response" {
		errs = append(errs, validate.ValidationError{
			Rule: "S-45", File: file, Func: funcName, SeqIdx: seqIdx, Level: "ERROR",
			Message: "@subscribe cannot use @response",
		})
	}
	for _, arg := range seq.Args {
		if arg.Source == "request" {
			errs = append(errs, validate.ValidationError{
				Rule: "S-40", File: file, Func: funcName, SeqIdx: seqIdx, Level: "ERROR",
				Message: "@subscribe cannot use request",
			})
		}
		if arg.Source == "query" {
			errs = append(errs, validate.ValidationError{
				Rule: "S-41", File: file, Func: funcName, SeqIdx: seqIdx, Level: "ERROR",
				Message: "@subscribe cannot use query",
			})
		}
	}
	return errs
}
