//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkPolicy — SSaC @auth ↔ Rego allow (action:resource) 교차 검증
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func checkPolicy(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if len(fs.Policies) == 0 {
		return nil
	}
	graph := toulmin.NewGraph("policy")
	graph.Rule(rule.PairMatch).With(&rule.PairMatchSpec{
		BaseSpec:  rule.BaseSpec{Rule: "X-28", Level: "WARNING", Message: "SSaC @auth pair has no matching Rego allow rule"},
		LookupKey: "Policy.auth",
	})

	var errs []CrossError
	for pair := range g.Pairs["SSaC.auth"] {
		ctx := toulmin.NewContext()
		ctx.Set("ground", g)
		ctx.Set("claim", pair)
		results, _ := graph.Evaluate(ctx)
		errs = append(errs, toErrors(results, pair)...)
	}
	return errs
}
