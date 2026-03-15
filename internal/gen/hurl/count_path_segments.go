//ff:func feature=gen-hurl type=util
//ff:what Counts the number of path segments in a URL path.
package hurl

import "strings"

func countPathSegments(path string) int {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	return len(parts)
}
