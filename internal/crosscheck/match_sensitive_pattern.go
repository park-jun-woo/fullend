//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=sensitive
//ff:what 컬럼 이름이 민감 패턴에 매치되면 해당 패턴을 반환
package crosscheck

import "strings"

func matchSensitivePattern(colName string) string {
	lower := strings.ToLower(colName)
	for _, p := range sensitivePatterns {
		if strings.Contains(lower, p) {
			return p
		}
	}
	return ""
}
