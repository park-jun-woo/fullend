//ff:func feature=contract type=util control=sequence
//ff:what 코멘트 라인이 fullend 디렉티브인지 확인한다
package contract

import "strings"

// IsDirective checks if a comment line is a fullend directive.
func IsDirective(comment string) bool {
	s := strings.TrimSpace(comment)
	return strings.HasPrefix(s, "//fullend:") || strings.HasPrefix(s, "// fullend:")
}
