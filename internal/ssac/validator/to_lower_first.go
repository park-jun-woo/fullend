//ff:func feature=ssac-validate type=util control=sequence
//ff:what 첫 글자를 소문자로 변환한다
package validator

import "strings"

// toLowerFirst는 첫 글자를 소문자로 변환한다.
func toLowerFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}
