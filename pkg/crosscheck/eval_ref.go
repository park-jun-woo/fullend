//ff:func feature=crosscheck type=util control=sequence
//ff:what evalRef — 기존 Graph로 단일 참조를 평가하여 CrossError 반환
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func evalRef(graph *toulmin.Graph, g *rule.Ground, claim, context string) []CrossError {
	ctx := toulmin.NewContext()
	ctx.Set("ground", g)
	ctx.Set("claim", claim)
	results, _ := graph.Evaluate(ctx)
	return toErrors(results, context)
}
