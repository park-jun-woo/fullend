//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkPolicyReverse — Rego allow → SSaC @auth 역방향 매칭 (X-29)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func checkPolicyReverse(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if len(fs.ParsedPolicies) == 0 {
		return nil
	}
	graph := toulmin.NewGraph("policy-reverse")
	graph.Rule(rule.PairMatch).With(&rule.PairMatchSpec{
		BaseSpec:  rule.BaseSpec{Rule: "X-29", Level: "WARNING", Message: "Rego allow rule has no matching SSaC @auth"},
		LookupKey: "SSaC.auth",
	})
	var errs []CrossError
	for pair := range g.Pairs["Policy.auth"] {
		errs = append(errs, evalRef(graph, g, pair, pair)...)
	}
	return errs
}
