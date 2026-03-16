//ff:func feature=gen-hurl type=util control=iteration dimension=3
//ff:what Role mapping — extracts operation -> required role from OPA policies.
package hurl

import "github.com/geul-org/fullend/internal/policy"

// buildOperationRoleMap extracts operation -> required role from OPA policies.
func buildOperationRoleMap(policies []*policy.Policy) map[string]string {
	roleMap := make(map[string]string)
	for _, p := range policies {
		for _, rule := range p.Rules {
			if rule.UsesRole && rule.RoleValue != "" {
				for _, action := range rule.Actions {
					roleMap[action] = rule.RoleValue
				}
			}
		}
	}
	return roleMap
}
