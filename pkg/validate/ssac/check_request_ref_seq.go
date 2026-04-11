//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what checkRequestRefSeq — 단일 시퀀스에서 request.field OpenAPI 스키마 존재 검증 (S-50 내부)
package ssac

import (
	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func checkRequestRefSeq(graph *toulmin.Graph, ground *rule.Ground, file, funcName string, seqIdx int, seq parsessac.Sequence) []validate.ValidationError {
	var errs []validate.ValidationError
	for _, arg := range seq.Args {
		if arg.Source != "request" || arg.Field == "" {
			continue
		}
		ctx := toulmin.NewContext()
		ctx.Set("ground", ground)
		ctx.Set("claim", arg.Field)
		results, _ := graph.Evaluate(ctx)
		errs = append(errs, toValidationErrors(results, file, funcName, seqIdx)...)
	}
	return errs
}
