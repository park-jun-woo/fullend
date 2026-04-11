//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what validateNameFormats — 이름 형식 검증 (S-26, S-46, S-47)
package ssac

import (
	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func validateNameFormats(fn parsessac.ServiceFunc, ground *rule.Ground) []validate.ValidationError {
	dotGraph := toulmin.NewGraph("name-dot-method")
	dotGraph.Rule(rule.NameFormat).With(&rule.NameFormatSpec{
		BaseSpec: rule.BaseSpec{Rule: "S-26", Level: "ERROR", Message: "Model must be Model.Method format"},
		Pattern:  "dot-method",
	})

	var errs []validate.ValidationError
	for i, seq := range fn.Sequences {
		if seq.Type == "get" || seq.Type == "post" || seq.Type == "put" || seq.Type == "delete" {
			ctx := toulmin.NewContext()
			ctx.Set("ground", ground)
			ctx.Set("claim", seq.Model)
			results, _ := dotGraph.Evaluate(ctx)
			errs = append(errs, toValidationErrors(results, fn.FileName, fn.Name, i)...)
		}
	}
	return errs
}
