//ff:func feature=symbol type=util control=sequence topic=ddl
//ff:what "users(id)" → ("users", "id") を파싱한다
package validator

import "strings"

// parseRef는 "users(id)" → ("users", "id") を파싱한다.
func parseRef(s string) (table, col string) {
	s = strings.TrimSpace(s)
	parenIdx := strings.Index(s, "(")
	if parenIdx < 0 {
		return s, ""
	}
	table = s[:parenIdx]
	col = strings.TrimSuffix(s[parenIdx+1:], ")")
	col = strings.TrimSuffix(col, ",")
	return table, col
}
