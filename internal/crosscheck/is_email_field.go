//ff:func feature=crosscheck type=util control=sequence topic=openapi-ddl
//ff:what 필드명이 email 패턴인지 확인
package crosscheck

import "strings"

func isEmailField(name string) bool {
	lower := strings.ToLower(name)
	return lower == "email" || strings.HasSuffix(lower, "email") || strings.HasPrefix(lower, "email")
}
