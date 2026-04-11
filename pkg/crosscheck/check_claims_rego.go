//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkClaimsRego — Rego input.claims → Config claims 존재 검증 (X-53)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func checkClaimsRego(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if len(fs.ParsedPolicies) == 0 {
		return nil
	}
	graph := toulmin.NewGraph("claims-rego")
	graph.Rule(rule.RefExists).With(&rule.RefExistsSpec{
		BaseSpec:  rule.BaseSpec{Rule: "X-53", Level: "ERROR", Message: "Rego input.claims reference not in fullend.yaml claims"},
		LookupKey: "Config.claims.keys",
	})

	var errs []CrossError
	for ref := range g.Lookup["Rego.claims"] {
		errs = append(errs, evalRef(graph, g, ref, ref)...)
	}
	return errs
}
