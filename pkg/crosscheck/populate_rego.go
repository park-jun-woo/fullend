//ff:func feature=crosscheck type=loader control=iteration dimension=1
//ff:what populateRego — Rego Policy에서 auth 쌍, ownership, claims, roles 추출
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func populateRego(g *rule.Ground, fs *fullend.Fullstack) {
	authPairs := make(rule.StringSet)
	claimsRefs := make(rule.StringSet)
	regoRoles := make(rule.StringSet)

	for _, p := range fs.ParsedPolicies {
		populateRegoPolicy(p, authPairs, claimsRefs, regoRoles)
	}
	g.Pairs["Policy.auth"] = authPairs
	g.Lookup["Rego.claims"] = claimsRefs
	g.Lookup["Rego.roles"] = regoRoles
}
