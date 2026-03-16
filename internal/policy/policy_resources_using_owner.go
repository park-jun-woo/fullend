//ff:func feature=policy type=util control=iteration dimension=1
//ff:what input.resource_owner를 참조하는 리소스 목록을 반환한다
package policy

// ResourcesUsingOwner returns resources that reference input.resource_owner in allow rules.
func (p *Policy) ResourcesUsingOwner() []string {
	seen := make(map[string]bool)
	for _, r := range p.Rules {
		if r.UsesOwner {
			seen[r.Resource] = true
		}
	}
	var result []string
	for res := range seen {
		result = append(result, res)
	}
	return result
}
