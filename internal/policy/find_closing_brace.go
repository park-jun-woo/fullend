//ff:func feature=policy type=util control=iteration dimension=1 topic=policy-check
//ff:what 중괄호 깊이를 추적하여 allow 블록의 닫는 중괄호 인덱스를 찾는다
package policy

// findClosingBrace finds the index of the closing brace that matches the opening
// of an allow block, accounting for nested braces (e.g., action sets).
func findClosingBrace(s string) int {
	depth := 1
	for i, c := range s {
		if c == '{' {
			depth++
		}
		if c == '}' {
			depth--
		}
		if depth == 0 {
			return i
		}
	}
	return -1
}
