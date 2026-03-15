//ff:func feature=gen-hurl type=util
//ff:what Returns sorted unique roles needed by operations.
package hurl

import "sort"

// collectRequiredRoles returns sorted unique roles needed by operations.
func collectRequiredRoles(roleMap map[string]string) []string {
	seen := make(map[string]bool)
	for _, role := range roleMap {
		seen[role] = true
	}
	var roles []string
	for role := range seen {
		roles = append(roles, role)
	}
	sort.Strings(roles)
	return roles
}
