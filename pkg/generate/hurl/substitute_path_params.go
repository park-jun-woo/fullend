//ff:func feature=gen-hurl type=util control=iteration dimension=1
//ff:what Replaces {ParamName} with {{captured_var}} using captured variables.
package hurl

import "strings"

// substitutePathParams replaces {ParamName} with {{captured_var}} using captured variables.
// For {ID}, it looks at the preceding path segment to find the matching capture (e.g. /gigs/{ID} -> gig_id).
func substitutePathParams(path string, captures map[string]bool) string {
	segments := strings.Split(path, "/")
	for i, seg := range segments {
		if !strings.HasPrefix(seg, "{") || !strings.HasSuffix(seg, "}") {
			continue
		}
		paramName := seg[1 : len(seg)-1]
		snakeParam := pascalToSnakeHurl(paramName)

		// First, check if the plain snake param exists in captures.
		if captures[snakeParam] {
			segments[i] = "{{" + snakeParam + "}}"
			continue
		}

		// Derive from preceding segment: /gigs/{ID} -> gig_id.
		derivedVar := ""
		if i > 0 {
			derivedVar = strings.TrimSuffix(segments[i-1], "s") + "_" + snakeParam
		}
		if derivedVar != "" && captures[derivedVar] {
			segments[i] = "{{" + derivedVar + "}}"
			continue
		}

		// Fallback: use plain snake param.
		segments[i] = "{{" + snakeParam + "}}"
	}
	return strings.Join(segments, "/")
}
