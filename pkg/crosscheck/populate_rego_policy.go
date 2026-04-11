//ff:func feature=crosscheck type=loader control=iteration dimension=1
//ff:what populateRegoPolicy — 단일 Policy에서 auth 쌍, claims, roles 추출
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/rego"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func populateRegoPolicy(p rego.Policy, authPairs, claimsRefs, regoRoles rule.StringSet) {
	for _, r := range p.Rules {
		for _, action := range r.Actions {
			authPairs[action+":"+r.Resource] = true
		}
		if r.UsesRole && r.RoleValue != "" {
			regoRoles[r.RoleValue] = true
		}
	}
	for _, ref := range p.ClaimsRefs {
		claimsRefs[ref] = true
	}
}
