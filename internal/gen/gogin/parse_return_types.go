//ff:func feature=gen-gogin type=util control=iteration dimension=1
//ff:what extracts individual return types from a return signature string

package gogin

import "strings"

// parseReturnTypes extracts individual return types from a return signature string.
// e.g. "(*Course, error)" → ["*Course", "error"]
func parseReturnTypes(sig string) []string {
	s := strings.TrimSpace(sig)
	s = strings.TrimPrefix(s, "(")
	s = strings.TrimSuffix(s, ")")
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
