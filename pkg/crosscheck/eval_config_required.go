//ff:func feature=crosscheck type=rule control=sequence
//ff:what evalConfigRequired — 단일 ConfigRequired 규칙을 평가하여 CrossError 반환
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func evalConfigRequired(g *rule.Ground, ruleID, configKey, message string) []CrossError {
	graph := toulmin.NewGraph("config-" + ruleID)
	graph.Rule(rule.ConfigRequired).With(&rule.ConfigRequiredSpec{
		BaseSpec:  rule.BaseSpec{Rule: ruleID, Level: "ERROR", Message: message},
		ConfigKey: configKey,
	})
	ctx := toulmin.NewContext()
	ctx.Set("ground", g)
	ctx.Set("claim", nil)
	results, _ := graph.Evaluate(ctx)
	return toErrors(results, "fullend.yaml")
}
