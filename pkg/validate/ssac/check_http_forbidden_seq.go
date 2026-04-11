//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what checkHTTPForbiddenSeq — 단일 시퀀스에서 message 사용 여부 검사 (S-44 내부)
package ssac

import (
	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func checkHTTPForbiddenSeq(graph *toulmin.Graph, ground *rule.Ground, file, funcName string, seqIdx int, seq parsessac.Sequence) []validate.ValidationError {
	var errs []validate.ValidationError
	for _, arg := range seq.Args {
		if arg.Source == "message" {
			ctx := toulmin.NewContext()
			ctx.Set("ground", ground)
			ctx.Set("claim", arg.Source)
			results, _ := graph.Evaluate(ctx)
			errs = append(errs, toValidationErrors(results, file, funcName, seqIdx)...)
		}
	}
	return errs
}
