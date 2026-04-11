//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what evalCallTypeArgs — @call의 각 arg에 대해 TypeMatch 평가
package ssac

import (
	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func evalCallTypeArgs(graph *toulmin.Graph, ground *rule.Ground, fn parsessac.ServiceFunc, seqIdx int, seq parsessac.Sequence, callFunc string) []validate.ValidationError {
	var errs []validate.ValidationError
	for _, arg := range seq.Args {
		if arg.Source == "" || arg.Field == "" {
			continue
		}
		sourceType := ground.Types["SSaC.var."+fn.Name+"."+arg.Source]
		if sourceType == "" {
			continue
		}
		ctx := toulmin.NewContext()
		ctx.Set("ground", ground)
		ctx.Set("claim", &rule.TypeClaim{Name: arg.Field, SourceType: sourceType})
		results, _ := graph.Evaluate(ctx)
		errs = append(errs, toValidationErrors(results, fn.FileName, fn.Name, seqIdx)...)
	}
	return errs
}
