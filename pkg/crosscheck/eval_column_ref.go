//ff:func feature=crosscheck type=util control=sequence
//ff:what evalColumnRef — 테이블별 DDL.column 참조 검증
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func evalColumnRef(g *rule.Ground, table, col, ruleID, context string) []CrossError {
	graph := toulmin.NewGraph("col-" + ruleID)
	graph.Rule(rule.RefExists).With(&rule.RefExistsSpec{
		BaseSpec:  rule.BaseSpec{Rule: ruleID, Level: "ERROR", Message: "column not found in DDL table " + table},
		LookupKey: "DDL.column." + table,
	})
	ctx := toulmin.NewContext()
	ctx.Set("ground", g)
	ctx.Set("claim", col)
	results, _ := graph.Evaluate(ctx)
	return toErrors(results, context)
}
