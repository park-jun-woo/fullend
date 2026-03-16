//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=config-check
//ff:what Rego 정책에서 input.claims 참조를 수집
package crosscheck

import "github.com/geul-org/fullend/internal/policy"

func collectRegoClaimsRefs(policies []*policy.Policy) map[string]string {
	regoRefs := make(map[string]string)
	for _, p := range policies {
		collectPolicyClaimsRefs(p, regoRefs)
	}
	return regoRefs
}
