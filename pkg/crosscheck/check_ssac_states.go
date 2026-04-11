//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkSSaCStates — SSaC @state → diagram 존재 검증
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func checkSSaCStates(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if len(fs.StateDiagrams) == 0 {
		return nil
	}
	graph := toulmin.NewGraph("ssac-states")
	graph.Rule(rule.RefExists).With(&rule.RefExistsSpec{
		BaseSpec:  rule.BaseSpec{Rule: "X-24", Level: "ERROR", Message: "SSaC @state references non-existent diagram"},
		LookupKey: "States.diagram",
	})

	var errs []CrossError
	for _, fn := range fs.ServiceFuncs {
		for _, ref := range collectStateRefs(fn.Sequences, fn.Name) {
			ctx := toulmin.NewContext()
			ctx.Set("ground", g)
			ctx.Set("claim", ref.diagramID)
			results, _ := graph.Evaluate(ctx)
			errs = append(errs, toErrors(results, ref.funcName)...)
		}
	}
	return errs
}
