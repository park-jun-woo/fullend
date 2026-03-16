//ff:func feature=ssac-gen type=util control=sequence topic=string-convert
//ff:what Go 컨벤션에 맞게 첫 글자를 대문자로 변환 (이니셜리즘 포함)
package generator

import "strings"

// ucFirst는 Go 컨벤션에 맞게 첫 글자를 대문자로 변환한다.
// 이니셜리즘이면 전부 대문자: "id" -> "ID", "url" -> "URL"
func ucFirst(s string) string {
	if s == "" {
		return s
	}
	if commonInitialisms[strings.ToUpper(s)] {
		return strings.ToUpper(s)
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
