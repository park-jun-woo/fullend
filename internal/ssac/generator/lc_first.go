//ff:func feature=ssac-gen type=util control=sequence topic=string-convert
//ff:what Go 컨벤션에 맞게 첫 단어를 소문자로 변환
package generator

import "strings"

// lcFirst는 Go 컨벤션에 맞게 첫 "단어"를 소문자로 변환한다.
// "ID" -> "id", "CourseID" -> "courseID", "HTTPClient" -> "httpClient"
func lcFirst(s string) string {
	if s == "" {
		return s
	}
	upper := countLeadingUpper(s)
	if upper == 0 {
		return s
	}
	if upper == 1 {
		return strings.ToLower(s[:1]) + s[1:]
	}
	if upper == len(s) {
		return strings.ToLower(s)
	}
	return strings.ToLower(s[:upper-1]) + s[upper-1:]
}
