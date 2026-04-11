//ff:func feature=crosscheck type=util control=sequence
//ff:what evalSchemaMatch — source 필드가 target 스키마에 존재하는지 평가
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func evalSchemaMatch(g *rule.Ground, source, target []string, ruleID, context string) []CrossError {
	localG := shallowCopyGround(g)
	localG.Schemas["_target"] = target
	graph := toulmin.NewGraph("schema-" + ruleID)
	graph.Rule(rule.SchemaMatch).With(&rule.SchemaMatchSpec{
		BaseSpec:  rule.BaseSpec{Rule: ruleID, Level: "ERROR", Message: "field not in target schema"},
		LookupKey: "_target",
	})
	ctx := toulmin.NewContext()
	ctx.Set("ground", localG)
	ctx.Set("claim", source)
	results, _ := graph.Evaluate(ctx)
	return toErrors(results, context)
}
