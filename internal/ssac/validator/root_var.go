//ff:func feature=ssac-validate type=util control=sequence
//ff:what "project.OwnerEmail" → "project" 루트 변수명을 추출한다
package validator

import "strings"

// rootVar는 "project.OwnerEmail" → "project"
func rootVar(s string) string {
	if idx := strings.Index(s, "."); idx >= 0 {
		return s[:idx]
	}
	return s
}
