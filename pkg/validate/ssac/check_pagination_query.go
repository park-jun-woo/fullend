//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what checkPaginationQuery — pagination @get에 {Query: query} 인자 존재 검증 (S-52)
package ssac

import (
	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

func checkPaginationQuery(file, funcName string, seqIdx int, seq parsessac.Sequence) []validate.ValidationError {
	for _, arg := range seq.Args {
		if arg.Source == "query" {
			return nil
		}
	}
	if seq.Inputs["Query"] != "" {
		return nil
	}
	return []validate.ValidationError{{
		Rule: "S-52", File: file, Func: funcName, SeqIdx: seqIdx, Level: "ERROR",
		Message: "pagination @get requires {Query: query} argument",
	}}
}
