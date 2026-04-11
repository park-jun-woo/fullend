//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what checkBindNotFound — fetch block의 bind가 response에도 custom.ts에도 없는지 검증
package stml

import (
	parsestml "github.com/park-jun-woo/fullend/pkg/parser/stml"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func checkBindNotFound(fb parsestml.FetchBlock, file string, ground *rule.Ground) []validate.ValidationError {
	respKey := "OpenAPI.response.resolved." + fb.OperationID
	graph := toulmin.NewGraph("stml-bind-check-" + fb.OperationID)
	w := graph.Rule(rule.RefExists).With(&rule.RefExistsSpec{
		BaseSpec:  rule.BaseSpec{Rule: "TM-8", Level: "ERROR", Message: "data-bind field not found in response or custom.ts"},
		LookupKey: respKey,
	})
	d := graph.Except(rule.IsCustomTS)
	d.Attacks(w)

	var errs []validate.ValidationError
	for _, b := range fb.Binds {
		ctx := toulmin.NewContext()
		ctx.Set("ground", ground)
		ctx.Set("claim", b.Name)
		results, _ := graph.Evaluate(ctx)
		errs = append(errs, toSTMLErrors(results, file, fb.OperationID)...)
	}
	return errs
}
