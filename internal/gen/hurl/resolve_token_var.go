//ff:func feature=gen-hurl type=util
//ff:what Picks the correct token variable for an operation based on role mapping.
package hurl

import (
	"sort"
	"strings"
)

// resolveTokenVar picks the correct token variable for an operation.
// If role-specific tokens exist (token_client, token_freelancer), uses the one matching
// the operation's required role. Falls back to the first role token alphabetically
// (deterministic) for owner-based operations without explicit role requirements.
func resolveTokenVar(operationID string, roleMap map[string]string, captures map[string]bool) string {
	if role, ok := roleMap[operationID]; ok {
		roleToken := "token_" + role
		if captures[roleToken] {
			return roleToken
		}
	}
	// Fallback: plain "token" (single-role mode).
	if captures["token"] {
		return "token"
	}
	// Pick the first role token alphabetically (deterministic).
	var tokens []string
	for k := range captures {
		if strings.HasPrefix(k, "token_") {
			tokens = append(tokens, k)
		}
	}
	if len(tokens) > 0 {
		sort.Strings(tokens)
		return tokens[0]
	}
	return ""
}
