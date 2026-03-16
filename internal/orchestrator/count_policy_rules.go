//ff:func feature=orchestrator type=util control=iteration dimension=1
//ff:what countPolicyRules counts total rules across all parsed policies.

package orchestrator

import (
	"github.com/geul-org/fullend/internal/policy"
)

func countPolicyRules(policies []*policy.Policy) int {
	total := 0
	for _, p := range policies {
		total += len(p.Rules)
	}
	return total
}
