//ff:func feature=crosscheck type=util control=sequence topic=policy-check
//ff:what extractQuotedLiteral — "..." 형식 문자열에서 따옴표 제거하고 content 반환

package crosscheck

import "strings"

func extractQuotedLiteral(v string) (string, bool) {
	v = strings.TrimSpace(v)
	if len(v) >= 2 && v[0] == '"' && v[len(v)-1] == '"' {
		return v[1 : len(v)-1], true
	}
	return "", false
}
