//ff:func feature=crosscheck type=util control=sequence topic=func-check
//ff:what rootVar — "x.y.z" → "x" (첫 점 앞의 변수명)

package crosscheck

import "strings"

func rootVar(target string) string {
	if i := strings.Index(target, "."); i >= 0 {
		return target[:i]
	}
	return target
}
