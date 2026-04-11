//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what validateAction — action ブロックの method/field 検証 (TM-2, TM-3, TM-5)
package stml

import (
	parsestml "github.com/park-jun-woo/fullend/pkg/parser/stml"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func validateAction(pages []parsestml.PageSpec, ground *rule.Ground) []validate.ValidationError {
	graph := toulmin.NewGraph("stml-action-ref")
	graph.Rule(rule.RefExists).With(&rule.RefExistsSpec{
		BaseSpec:  rule.BaseSpec{Rule: "TM-2", Level: "ERROR", Message: "data-action operationId not in OpenAPI"},
		LookupKey: "OpenAPI.operationId",
	})

	var errs []validate.ValidationError
	for _, page := range pages {
		for _, ab := range page.Actions {
			ctx := toulmin.NewContext()
			ctx.Set("ground", ground)
			ctx.Set("claim", ab.OperationID)
			results, _ := graph.Evaluate(ctx)
			errs = append(errs, toSTMLErrors(results, page.FileName, ab.OperationID)...)
		}
	}
	return errs
}
