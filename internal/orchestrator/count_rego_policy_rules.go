//ff:func feature=orchestrator type=util control=iteration dimension=1
//ff:what countRegoPolicyRules — pkg/parser/rego.Policy 슬라이스의 총 rule 수

package orchestrator

import "github.com/park-jun-woo/fullend/pkg/parser/rego"

func countRegoPolicyRules(policies []rego.Policy) int {
	total := 0
	for _, p := range policies {
		total += len(p.Rules)
	}
	return total
}
