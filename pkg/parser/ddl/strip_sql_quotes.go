//ff:func feature=manifest type=util control=sequence topic=ddl
//ff:what stripSQLQuotes — '...' 또는 "..." 문자열에서 양끝 따옴표 제거

package ddl

import "strings"

func stripSQLQuotes(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 {
		if (s[0] == '\'' && s[len(s)-1] == '\'') || (s[0] == '"' && s[len(s)-1] == '"') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
