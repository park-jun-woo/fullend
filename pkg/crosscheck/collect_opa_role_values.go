//ff:func feature=crosscheck type=util control=iteration dimension=2 topic=policy-check
//ff:what collectOPARoleValues — ParsedPolicies 에서 언급된 모든 role 값 + 제약없는 allow 보유 여부

package crosscheck

import "github.com/park-jun-woo/fullend/pkg/fullend"

// collectOPARoleValues collects role strings explicitly required by any allow rule,
// and reports whether any allow rule is unconstrained by role (UsesRole=false).
// Unconstrained allow means role restriction may be checked elsewhere.
func collectOPARoleValues(fs *fullend.Fullstack) (roles map[string]bool, hasUnconstrained bool) {
	roles = map[string]bool{}
	for _, p := range fs.ParsedPolicies {
		for _, r := range p.Rules {
			if r.UsesRole && r.RoleValue != "" {
				roles[r.RoleValue] = true
			} else {
				hasUnconstrained = true
			}
		}
	}
	return roles, hasUnconstrained
}
