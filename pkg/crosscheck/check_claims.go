//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkClaims — SSaC currentUser.Field → Config claims 존재 검증 (X-49)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func checkClaims(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	graph := toulmin.NewGraph("claims")
	graph.Rule(rule.RefExists).With(&rule.RefExistsSpec{
		BaseSpec:  rule.BaseSpec{Rule: "X-49", Level: "ERROR", Message: "currentUser field not in fullend.yaml claims"},
		LookupKey: "Config.claims",
	})

	var errs []CrossError
	for _, field := range collectCurrentUserFields(fs) {
		errs = append(errs, evalRef(graph, g, field, "currentUser."+field)...)
	}
	return errs
}
