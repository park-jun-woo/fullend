//ff:func feature=ssac-gen type=util control=sequence topic=string-convert
//ff:what dotted 변수명에서 루트 변수를 추출 (예: "a.b" -> "a")
package ssac

import "strings"

func rootVar(s string) string {
	if idx := strings.Index(s, "."); idx >= 0 {
		return s[:idx]
	}
	return s
}
