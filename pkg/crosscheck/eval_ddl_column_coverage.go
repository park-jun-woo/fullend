//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what evalDDLColumnCoverage — DDL 컬럼이 OpenAPI response에 포함되는지 평가 (X-10)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func evalDDLColumnCoverage(g *rule.Ground, tableName string, responseFields []string) []CrossError {
	respSet := make(rule.StringSet, len(responseFields))
	for _, f := range responseFields {
		respSet[f] = true
	}
	localG := shallowCopyGround(g)
	localG.Lookup["_response"] = respSet

	graph := toulmin.NewGraph("ddl-oapi-coverage")
	w := graph.Rule(rule.CoverageCheck).With(&rule.CoverageCheckSpec{
		BaseSpec:  rule.BaseSpec{Rule: "X-10", Level: "WARNING", Message: "DDL column not in OpenAPI response"},
		LookupKey: "_response",
	})
	d := graph.Except(rule.IsSensitiveCol)
	d.Attacks(w)

	var errs []CrossError
	for col := range g.Lookup["DDL.column."+tableName] {
		ctx := toulmin.NewContext()
		ctx.Set("ground", localG)
		ctx.Set("claim", col)
		results, _ := graph.Evaluate(ctx)
		errs = append(errs, toErrors(results, tableName+"."+col)...)
	}
	return errs
}
