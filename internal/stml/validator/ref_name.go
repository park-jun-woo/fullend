//ff:func feature=stml-validate type=util control=sequence
//ff:what $ref 문자열에서 스키마 이름만 추출
package validator

import "strings"

func refName(ref string) string {
	parts := strings.Split(ref, "/")
	return parts[len(parts)-1]
}
