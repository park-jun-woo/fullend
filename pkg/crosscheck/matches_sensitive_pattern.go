//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what matchesSensitivePattern — 컬럼명이 민감 패턴에 해당하는지 확인
package crosscheck

import "strings"

var sensitivePatterns = []string{"password", "secret", "hash", "token"}

func matchesSensitivePattern(col string) bool {
	lower := strings.ToLower(col)
	for _, p := range sensitivePatterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}
