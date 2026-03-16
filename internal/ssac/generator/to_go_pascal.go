//ff:func feature=ssac-gen type=util control=sequence topic=string-convert
//ff:what 문자열의 첫 글자를 대문자로 변환하는 경량 래퍼
package generator

import "strings"

// toGoPascal is a local wrapper to convert first letter to uppercase.
func toGoPascal(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
