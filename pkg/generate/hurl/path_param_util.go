//ff:func feature=gen-hurl type=util control=iteration dimension=1
//ff:what Path parameter handling — checks if all path parameters can be resolved from captured variables.
package hurl

import "strings"

// canResolvePathParams checks if all path parameters in a path can be resolved
// from captured variables. Returns false if any param would be undefined at runtime.
func canResolvePathParams(path string, captures map[string]bool) bool {
	segments := strings.Split(path, "/")
	for i, seg := range segments {
		if !strings.HasPrefix(seg, "{") || !strings.HasSuffix(seg, "}") {
			continue
		}
		paramName := seg[1 : len(seg)-1]
		snakeParam := pascalToSnakeHurl(paramName)

		// Check plain snake param.
		if captures[snakeParam] {
			continue
		}

		// Check derived from preceding segment: /gigs/{ID} -> gig_id.
		if i > 0 && captures[strings.TrimSuffix(segments[i-1], "s")+"_"+snakeParam] {
			continue
		}

		// Unresolvable param found.
		return false
	}
	return true
}
