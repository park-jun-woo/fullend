//ff:func feature=gen-hurl type=util control=iteration dimension=1
//ff:what Finds a captured variable that matches an FK field name via direct or prefix match.
package hurl

import "strings"

// findMatchingCapture finds a captured variable that matches an FK field name.
// e.g. fieldName="org_id", captures has "organization_id" -> returns "organization_id".
// Tries direct match first, then prefix match (org -> organization).
func findMatchingCapture(fieldName string, captures map[string]bool) string {
	// Direct match.
	if captures[fieldName] {
		return fieldName
	}
	// Prefix match: org_id -> prefix "org", find "*_id" capture where prefix matches.
	if !strings.HasSuffix(fieldName, "_id") {
		return ""
	}
	prefix := strings.TrimSuffix(fieldName, "_id")
	for cap := range captures {
		if !strings.HasSuffix(cap, "_id") {
			continue
		}
		if strings.HasPrefix(strings.TrimSuffix(cap, "_id"), prefix) {
			return cap
		}
	}
	return ""
}
