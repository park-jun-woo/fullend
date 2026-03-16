//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what Rego 정책에서 role 값을 수집
package crosscheck

import "github.com/geul-org/fullend/internal/policy"

func collectRegoRoles(policies []*policy.Policy) map[string]string {
	regoRoles := make(map[string]string)
	for _, p := range policies {
		collectPolicyRoles(p, regoRoles)
	}
	return regoRoles
}
