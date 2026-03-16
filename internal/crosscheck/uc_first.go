//ff:func feature=crosscheck type=util control=sequence
//ff:what 첫 글자를 대문자로 변환 (PascalCase)
package crosscheck

import "strings"

// ucFirst converts first letter to uppercase (PascalCase).
func ucFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
