//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkStates — States transition ↔ SSaC/OpenAPI 교차 검증
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func checkStates(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if len(fs.StateDiagrams) == 0 {
		return nil
	}
	graph := toulmin.NewGraph("states")
	graph.Rule(rule.RefExists).With(&rule.RefExistsSpec{
		BaseSpec:  rule.BaseSpec{Rule: "X-23", Level: "ERROR", Message: "States transition event has no matching SSaC function"},
		LookupKey: "SSaC.funcName",
	})

	var errs []CrossError
	for _, sd := range fs.StateDiagrams {
		for _, tr := range sd.Transitions {
			ctx := toulmin.NewContext()
			ctx.Set("ground", g)
			ctx.Set("claim", tr.Event)
			results, _ := graph.Evaluate(ctx)
			errs = append(errs, toErrors(results, sd.ID+"."+tr.Event)...)
		}
	}
	return errs
}
