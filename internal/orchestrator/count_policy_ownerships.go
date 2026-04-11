//ff:func feature=orchestrator type=util control=iteration dimension=1
//ff:what countPolicyOwnerships — policy에서 총 ownership 매핑 수 집계
package orchestrator

import "github.com/park-jun-woo/fullend/internal/policy"

func countPolicyOwnerships(policies []*policy.Policy) int {
	total := 0
	for _, p := range policies {
		total += len(p.Ownerships)
	}
	return total
}
