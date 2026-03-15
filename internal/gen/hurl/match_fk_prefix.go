//ff:func feature=gen-hurl type=util
//ff:what Returns true if the resource name starts with any FK prefix.
package hurl

import "strings"

// matchFKPrefix returns true if the resource name starts with any FK prefix.
// e.g. resource="organizations", prefix="org" -> true (organizations starts with org).
func matchFKPrefix(resource string, fkPrefixes []string) bool {
	for _, prefix := range fkPrefixes {
		if strings.HasPrefix(resource, prefix) {
			return true
		}
	}
	return false
}
