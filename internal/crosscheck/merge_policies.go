//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=policy-check
//ff:what 모든 정책의 action-resource 쌍과 ownership 정보를 병합
package crosscheck

import "github.com/geul-org/fullend/internal/policy"

func mergePolicies(policies []*policy.Policy) (map[[2]string]bool, map[string]bool, []policy.OwnershipMapping) {
	allPairs := make(map[[2]string]bool)
	ownerResources := make(map[string]bool)
	var allOwnerships []policy.OwnershipMapping

	for _, p := range policies {
		for _, pair := range p.ActionResourcePairs() {
			allPairs[pair] = true
		}
		for _, res := range p.ResourcesUsingOwner() {
			ownerResources[res] = true
		}
		allOwnerships = append(allOwnerships, p.Ownerships...)
	}

	return allPairs, ownerResources, allOwnerships
}
