//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what checkActionFieldRefs — action block의 data-field가 OpenAPI request에 있는지 검증
package stml

import (
	parsestml "github.com/park-jun-woo/fullend/pkg/parser/stml"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func checkActionFieldRefs(ab parsestml.ActionBlock, file string, ground *rule.Ground) []validate.ValidationError {
	reqKey := "OpenAPI.request." + ab.OperationID
	if _, ok := ground.Lookup[reqKey]; !ok {
		return nil
	}
	graph := toulmin.NewGraph("stml-field-" + ab.OperationID)
	graph.Rule(rule.RefExists).With(&rule.RefExistsSpec{
		BaseSpec:  rule.BaseSpec{Rule: "TM-5", Level: "ERROR", Message: "data-field not in OpenAPI request schema"},
		LookupKey: reqKey,
	})
	var errs []validate.ValidationError
	for _, f := range ab.Fields {
		ctx := toulmin.NewContext()
		ctx.Set("ground", ground)
		ctx.Set("claim", f.Name)
		results, _ := graph.Evaluate(ctx)
		errs = append(errs, toSTMLErrors(results, file, ab.OperationID)...)
	}
	return errs
}
