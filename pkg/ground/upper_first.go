//ff:func feature=rule type=util control=sequence
//ff:what upperFirst — 문자열의 첫 글자를 대문자로 변환
package ground

import "strings"

func upperFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
