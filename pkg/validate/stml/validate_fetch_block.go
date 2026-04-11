//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what validateFetchBlock — 단일 fetch 블록의 bind/param 참조 검증
package stml

import (
	parsestml "github.com/park-jun-woo/fullend/pkg/parser/stml"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func validateFetchBlock(fb parsestml.FetchBlock, file string, ground *rule.Ground) []validate.ValidationError {
	var errs []validate.ValidationError

	// TM-4: data-param → OpenAPI param
	paramGraph := toulmin.NewGraph("stml-param-" + fb.OperationID)
	paramGraph.Rule(rule.RefExists).With(&rule.RefExistsSpec{
		BaseSpec:  rule.BaseSpec{Rule: "TM-4", Level: "ERROR", Message: "data-param not in OpenAPI parameters"},
		LookupKey: "OpenAPI.param." + fb.OperationID,
	})
	for _, p := range fb.Params {
		ctx := toulmin.NewContext()
		ctx.Set("ground", ground)
		ctx.Set("claim", p.Name)
		results, _ := paramGraph.Evaluate(ctx)
		errs = append(errs, toSTMLErrors(results, file, fb.OperationID)...)
	}

	// TM-6: data-bind → OpenAPI response field
	bindGraph := toulmin.NewGraph("stml-bind-" + fb.OperationID)
	w := bindGraph.Rule(rule.RefExists).With(&rule.RefExistsSpec{
		BaseSpec:  rule.BaseSpec{Rule: "TM-6", Level: "ERROR", Message: "data-bind field not in OpenAPI response"},
		LookupKey: "OpenAPI.response.resolved." + fb.OperationID,
	})
	d := bindGraph.Except(rule.IsCustomTS)
	d.Attacks(w)
	for _, b := range fb.Binds {
		ctx := toulmin.NewContext()
		ctx.Set("ground", ground)
		ctx.Set("claim", b.Name)
		results, _ := bindGraph.Evaluate(ctx)
		errs = append(errs, toSTMLErrors(results, file, fb.OperationID)...)
	}

	// TM-10: data-sort column → OpenAPI x-sort allowed
	if fb.Sort != nil {
		errs = append(errs, validateSortRef(fb, file, ground)...)
	}

	// TM-11: data-filter column → OpenAPI x-filter allowed
	if len(fb.Filters) > 0 {
		errs = append(errs, validateFilterRef(fb, file, ground)...)
	}

	return errs
}
