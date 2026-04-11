//ff:func feature=rule type=rule control=sequence
//ff:what checkSeqForbidden — 단일 시퀀스의 result 변수/모델 이름에 대해 금지 참조 검증
package ssac

import (
	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func checkSeqForbidden(goGraph, reservedGraph, dotGraph *toulmin.Graph, ground *rule.Ground, file, funcName string, seqIdx int, seq parsessac.Sequence) []validate.ValidationError {
	var errs []validate.ValidationError
	// S-34: result variable vs Go reserved
	if seq.Result != nil && seq.Result.Var != "" {
		ctx := toulmin.NewContext()
		ctx.Set("ground", ground)
		ctx.Set("claim", seq.Result.Var)
		results, _ := goGraph.Evaluate(ctx)
		errs = append(errs, toValidationErrors(results, file, funcName, seqIdx)...)
		// S-33: result variable vs reserved source
		results2, _ := reservedGraph.Evaluate(ctx)
		errs = append(errs, toValidationErrors(results2, file, funcName, seqIdx)...)
	}
	// S-47: Model must not have package prefix for @get/@post/@put/@delete
	if (seq.Type == "get" || seq.Type == "post" || seq.Type == "put" || seq.Type == "delete") && seq.Model != "" {
		model := extractModelFromSeq(seq)
		if model != "" {
			ctx := toulmin.NewContext()
			ctx.Set("ground", ground)
			ctx.Set("claim", model)
			results, _ := dotGraph.Evaluate(ctx)
			errs = append(errs, toValidationErrors(results, file, funcName, seqIdx)...)
		}
	}
	return errs
}
