//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what validateResultReserved — result type이 Go 예약어인지 검증 (S-35)
package ssac

import (
	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func validateResultReserved(fn parsessac.ServiceFunc, ground *rule.Ground) []validate.ValidationError {
	graph := toulmin.NewGraph("result-reserved")
	graph.Rule(rule.ForbiddenRef).With(&rule.ForbiddenRefSpec{
		BaseSpec:  rule.BaseSpec{Rule: "S-35", Level: "ERROR", Message: "result type is Go reserved word"},
		LookupKey: "go.reserved",
	})
	var errs []validate.ValidationError
	for i, seq := range fn.Sequences {
		if seq.Result == nil || seq.Result.Type == "" {
			continue
		}
		typeName := stripTypePrefix(seq.Result.Type)
		ctx := toulmin.NewContext()
		ctx.Set("ground", ground)
		ctx.Set("claim", typeName)
		results, _ := graph.Evaluate(ctx)
		errs = append(errs, toValidationErrors(results, fn.FileName, fn.Name, i)...)
	}
	return errs
}
