//ff:func feature=policy type=parser control=iteration dimension=1
//ff:what extractAllowRules — Rego 소스에서 allow 블록들을 분리하고 AllowRule 추출
package rego

import "strings"

func extractAllowRules(content string, p *Policy) {
	normalized := strings.ReplaceAll(content, "\nallow {", "\nallow if {")
	parts := strings.Split(normalized, "allow if {")
	for i := 1; i < len(parts); i++ {
		endIdx := findClosingBrace(parts[i])
		if endIdx < 0 {
			continue
		}
		if rule, ok := processAllowBlock(parts[i][:endIdx]); ok {
			p.Rules = append(p.Rules, rule)
		}
	}
}
