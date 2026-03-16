//ff:func feature=gen-hurl type=util control=iteration
//ff:what Extracts the parent resource from a nested path.
package hurl

import "strings"

// findParentResource extracts the parent resource from a nested path.
// e.g. /gigs/{GigID}/proposals -> "gigs", /gigs -> "".
func findParentResource(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	for i := 1; i < len(parts); i++ {
		if strings.HasPrefix(parts[i], "{") && i > 0 {
			// Found a path param; the segment before it is a resource.
			// If there's more path after this param, parts[i-1] is the parent.
			if i+1 < len(parts) {
				return parts[i-1]
			}
		}
	}
	return ""
}
