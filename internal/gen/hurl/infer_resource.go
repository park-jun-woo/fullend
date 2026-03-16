//ff:func feature=gen-hurl type=util control=iteration dimension=1
//ff:what Infers the resource name from the first non-parameter path segment.
package hurl

import "strings"

func inferResource(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	// Find the first non-parameter segment.
	for _, p := range parts {
		if !strings.HasPrefix(p, "{") {
			return p
		}
	}
	return "other"
}
