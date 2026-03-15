//ff:func feature=symbol type=util
//ff:what 괄호 안 컬럼을 추출한다
package validator

import "strings"

// extractParenColumns는 "PRIMARY KEY (col1, col2)" 등에서 괄호 안 컬럼을 추출한다.
func extractParenColumns(line string) []string {
	parenIdx := strings.Index(line, "(")
	if parenIdx < 0 {
		return nil
	}
	inner := line[parenIdx+1:]
	inner = strings.TrimSuffix(strings.TrimSpace(inner), ",")
	inner = strings.TrimSuffix(inner, ");")
	inner = strings.TrimSuffix(inner, ")")
	var cols []string
	for _, c := range strings.Split(inner, ",") {
		c = strings.TrimSpace(c)
		if c != "" {
			cols = append(cols, c)
		}
	}
	return cols
}
