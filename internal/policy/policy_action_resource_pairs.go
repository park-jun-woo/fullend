//ff:func feature=policy type=util control=iteration dimension=1
//ff:what 정책의 모든 (action, resource) 쌍을 반환한다
package policy

// ActionResourcePairs returns all (action, resource) pairs from the policy.
func (p *Policy) ActionResourcePairs() [][2]string {
	var pairs [][2]string
	for _, r := range p.Rules {
		for _, a := range r.Actions {
			pairs = append(pairs, [2]string{a, r.Resource})
		}
	}
	return pairs
}
