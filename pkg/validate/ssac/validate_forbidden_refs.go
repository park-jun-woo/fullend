//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what validateForbiddenRefs — 금지 참조 검증 (S-31~S-35, S-42~S-44)
package ssac

import (
	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func validateForbiddenRefs(fn parsessac.ServiceFunc, ground *rule.Ground) []validate.ValidationError {
	graph := toulmin.NewGraph("forbidden-refs")
	graph.Rule(rule.ForbiddenRef).With(&rule.ForbiddenRefSpec{
		BaseSpec:  rule.BaseSpec{Rule: "S-34", Level: "ERROR", Message: "Go reserved word used as variable name"},
		LookupKey: "go.reserved",
	})

	var errs []validate.ValidationError
	for i, seq := range fn.Sequences {
		if seq.Result != nil && seq.Result.Var != "" {
			ctx := toulmin.NewContext()
			ctx.Set("ground", ground)
			ctx.Set("claim", seq.Result.Var)
			results, _ := graph.Evaluate(ctx)
			errs = append(errs, toValidationErrors(results, fn.FileName, fn.Name, i)...)
		}
	}
	return errs
}
