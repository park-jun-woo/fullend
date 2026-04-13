//ff:func feature=gen-hurl type=util control=iteration dimension=1
//ff:what Extracts the parent resource from a nested path.
package hurl

import "strings"

// findParentResource extracts the parent resource from a nested path.
// e.g. /gigs/{GigID}/proposals -> "gigs", /gigs -> "".
func findParentResource(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	for i := 1; i < len(parts); i++ {
		if strings.HasPrefix(parts[i], "{") && i+1 < len(parts) {
			return parts[i-1]
		}
	}
	return ""
}
