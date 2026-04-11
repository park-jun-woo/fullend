//ff:func feature=rule type=rule control=sequence
//ff:what evalModelRef — 단일 시퀀스의 Model 참조와 result 타입 검증
package ssac

import (
	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func evalModelRef(modelGraph, upperGraph *toulmin.Graph, ground *rule.Ground, file, funcName string, seqIdx int, seq parsessac.Sequence) []validate.ValidationError {
	var errs []validate.ValidationError
	model := extractModelFromSeq(seq)
	if model != "" && (seq.Type == "get" || seq.Type == "post" || seq.Type == "put" || seq.Type == "delete") {
		ctx := toulmin.NewContext()
		ctx.Set("ground", ground)
		ctx.Set("claim", model)
		results, _ := modelGraph.Evaluate(ctx)
		errs = append(errs, toValidationErrors(results, file, funcName, seqIdx)...)
	}
	if seq.Result != nil && seq.Result.Type != "" {
		ctx := toulmin.NewContext()
		ctx.Set("ground", ground)
		ctx.Set("claim", seq.Result.Type)
		results, _ := upperGraph.Evaluate(ctx)
		errs = append(errs, toValidationErrors(results, file, funcName, seqIdx)...)
	}
	return errs
}
