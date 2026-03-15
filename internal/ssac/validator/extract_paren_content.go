//ff:func feature=symbol type=util
//ff:what "(content)" 에서 content를 추출한다
package validator

import "strings"

// extractParenContent는 "(content)" 에서 content를 추출한다.
func extractParenContent(s string) string {
	open := strings.Index(s, "(")
	close := strings.Index(s, ")")
	if open < 0 || close < 0 || close <= open {
		return ""
	}
	return strings.TrimSpace(s[open+1 : close])
}
