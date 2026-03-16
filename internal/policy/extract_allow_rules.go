//ff:func feature=policy type=parser control=iteration dimension=1
//ff:what Rego 소스에서 allow 블록들을 분리하고 AllowRule을 추출한다
package policy

import "strings"

// extractAllowRules extracts (action, resource) pairs from allow blocks.
func extractAllowRules(content string, p *Policy) {
	// Split by "allow if {" or "allow {" blocks (Rego v1 and v0 syntax).
	// Normalize "allow {" → "allow if {" for uniform parsing.
	normalized := strings.ReplaceAll(content, "\nallow {", "\nallow if {")
	parts := strings.Split(normalized, "allow if {")
	for i := 1; i < len(parts); i++ {
		block := parts[i]
		// Find the closing brace that matches the allow block (depth-aware).
		endIdx := findClosingBrace(block)
		if endIdx < 0 {
			continue
		}
		block = block[:endIdx]

		if rule, ok := processAllowBlock(block); ok {
			p.Rules = append(p.Rules, rule)
		}
	}
}
