//ff:func feature=manifest type=util control=iteration dimension=1 topic=ddl
//ff:what splitCSVLiterals — ',' 구분, 단 '...' 안의 쉼표 보존

package ddl

import "strings"

// splitCSVLiterals splits a VALUES parens body by commas while respecting SQL
// single-quoted string literals (commas inside '...' are preserved).
func splitCSVLiterals(s string) []string {
	var out []string
	var cur strings.Builder
	inQuote := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '\'' {
			inQuote = !inQuote
			cur.WriteByte(c)
			continue
		}
		if c == ',' && !inQuote {
			out = append(out, strings.TrimSpace(cur.String()))
			cur.Reset()
			continue
		}
		cur.WriteByte(c)
	}
	if cur.Len() > 0 {
		out = append(out, strings.TrimSpace(cur.String()))
	}
	return out
}
