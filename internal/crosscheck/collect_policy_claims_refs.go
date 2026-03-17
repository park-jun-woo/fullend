//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=config-check
//ff:what 단일 Rego 정책에서 claims 참조를 수집
package crosscheck

import "github.com/park-jun-woo/fullend/internal/policy"

func collectPolicyClaimsRefs(p *policy.Policy, regoRefs map[string]string) {
	for _, ref := range p.ClaimsRefs {
		if _, exists := regoRefs[ref]; !exists {
			regoRefs[ref] = p.File
		}
	}
}
