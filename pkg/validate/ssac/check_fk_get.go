//ff:func feature=rule type=rule control=sequence
//ff:what checkFKGet — 단일 @get 시퀀스의 FK 참조 + @empty 가드 검증
package ssac

import (
	"strings"

	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

func checkFKGet(fn parsessac.ServiceFunc, seqIdx int, seq parsessac.Sequence, declared map[string]bool, varTypes map[string]string) []validate.ValidationError {
	if seq.Type != "get" || seq.Result == nil {
		return nil
	}
	if strings.HasPrefix(seq.Result.Type, "[]") || seq.Result.Wrapper != "" {
		return nil
	}
	getModel := extractModelFromSeq(seq)
	if !hasFKRefInArgs(seq, declared, varTypes, getModel) {
		return nil
	}
	if hasEmptyGuard(fn.Sequences[seqIdx+1:], seq.Result.Var) {
		return nil
	}
	return []validate.ValidationError{{
		Rule: "S-37", File: fn.FileName, Func: fn.Name, SeqIdx: seqIdx, Level: "WARNING",
		Message: seq.Result.Var + " — FK reference @get requires @empty guard",
	}}
}
