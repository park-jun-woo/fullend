//ff:func feature=contract type=util control=selection
//ff:what 코멘트 문자열에서 fullend 디렉티브 접두사를 제거한다
package contract

import (
	"fmt"
	"strings"
)

// stripDirectivePrefix removes the //fullend: or // fullend: prefix and returns the remainder.
func stripDirectivePrefix(comment string) (string, error) {
	s := strings.TrimSpace(comment)
	switch {
	case strings.HasPrefix(s, "//fullend:"):
		return strings.TrimPrefix(s, "//fullend:"), nil
	case strings.HasPrefix(s, "// fullend:"):
		return strings.TrimPrefix(s, "// fullend:"), nil
	default:
		return "", fmt.Errorf("not a fullend directive: %q", comment)
	}
}
