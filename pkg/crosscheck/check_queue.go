//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkQueue — @publish ↔ @subscribe 토픽 교차 검증
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func checkQueue(g *rule.Ground) []CrossError {
	if len(g.Pairs["SSaC.publish"]) == 0 && len(g.Pairs["SSaC.subscribe"]) == 0 {
		return nil
	}
	var errs []CrossError

	pubGraph := toulmin.NewGraph("queue-pub")
	pubGraph.Rule(rule.PairMatch).With(&rule.PairMatchSpec{
		BaseSpec:  rule.BaseSpec{Rule: "X-57", Level: "WARNING", Message: "@publish topic has no @subscribe handler"},
		LookupKey: "SSaC.subscribe",
	})
	for topic := range g.Pairs["SSaC.publish"] {
		ctx := toulmin.NewContext()
		ctx.Set("ground", g)
		ctx.Set("claim", topic)
		results, _ := pubGraph.Evaluate(ctx)
		errs = append(errs, toErrors(results, topic)...)
	}

	subGraph := toulmin.NewGraph("queue-sub")
	subGraph.Rule(rule.PairMatch).With(&rule.PairMatchSpec{
		BaseSpec:  rule.BaseSpec{Rule: "X-58", Level: "WARNING", Message: "@subscribe topic has no @publish source"},
		LookupKey: "SSaC.publish",
	})
	for topic := range g.Pairs["SSaC.subscribe"] {
		ctx := toulmin.NewContext()
		ctx.Set("ground", g)
		ctx.Set("claim", topic)
		results, _ := subGraph.Evaluate(ctx)
		errs = append(errs, toErrors(results, topic)...)
	}

	return errs
}
