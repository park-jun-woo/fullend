//ff:func feature=crosscheck type=util control=sequence topic=openapi-ddl
//ff:what 필드명이 password 패턴인지 확인
package crosscheck

import "strings"

func isPasswordField(name string) bool {
	lower := strings.ToLower(name)
	return lower == "password" || strings.HasSuffix(lower, "password") || strings.HasPrefix(lower, "password")
}
