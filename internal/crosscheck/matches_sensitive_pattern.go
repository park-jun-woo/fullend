//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what 컬럼명이 민감 정보 패턴과 일치하는지 확인
package crosscheck

import "strings"

// matchesSensitivePattern checks if a column name matches any sensitive pattern.
func matchesSensitivePattern(colName string) bool {
	lower := strings.ToLower(colName)
	for _, p := range sensitivePatterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}
