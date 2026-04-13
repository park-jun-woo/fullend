//ff:func feature=pkg-authz type=util control=iteration dimension=1
//ff:what toSnakeCase — CamelCase → snake_case (연속 대문자 acronym 보존)

package authz

import (
	"strings"
	"unicode"
)

func toSnakeCase(s string) string {
	var b strings.Builder
	runes := []rune(s)
	for i, r := range runes {
		if shouldInsertUnderscore(runes, i) {
			b.WriteByte('_')
		}
		b.WriteRune(unicode.ToLower(r))
	}
	return b.String()
}
