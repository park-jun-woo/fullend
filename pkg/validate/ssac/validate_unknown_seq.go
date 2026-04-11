//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what validateUnknownSeq — 알 수 없는 시퀀스 타입 검증 (S-25)
package ssac

import (
	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

var knownSeqTypes = map[string]bool{
	"get": true, "post": true, "put": true, "delete": true,
	"empty": true, "exists": true, "state": true, "auth": true,
	"call": true, "publish": true, "response": true,
}

func validateUnknownSeq(fn parsessac.ServiceFunc) []validate.ValidationError {
	var errs []validate.ValidationError
	for i, seq := range fn.Sequences {
		if !knownSeqTypes[seq.Type] {
			errs = append(errs, validate.ValidationError{
				Rule: "S-25", File: fn.FileName, Func: fn.Name, SeqIdx: i, Level: "ERROR",
				Message: "unknown sequence type: @" + seq.Type,
			})
		}
	}
	return errs
}
