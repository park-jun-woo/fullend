//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what validateFilterRef — data-filter 컬럼이 OpenAPI x-filter allowed에 있는지 검증 (TM-11)
package stml

import (
	parsestml "github.com/park-jun-woo/fullend/pkg/parser/stml"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func validateFilterRef(fb parsestml.FetchBlock, file string, ground *rule.Ground) []validate.ValidationError {
	graph := toulmin.NewGraph("stml-filter-" + fb.OperationID)
	graph.Rule(rule.RefExists).With(&rule.RefExistsSpec{
		BaseSpec:  rule.BaseSpec{Rule: "TM-11", Level: "ERROR", Message: "data-filter column not in OpenAPI x-filter allowed"},
		LookupKey: "OpenAPI.filter." + fb.OperationID,
	})
	var errs []validate.ValidationError
	for _, col := range fb.Filters {
		ctx := toulmin.NewContext()
		ctx.Set("ground", ground)
		ctx.Set("claim", col)
		results, _ := graph.Evaluate(ctx)
		errs = append(errs, toSTMLErrors(results, file, fb.OperationID)...)
	}
	return errs
}
