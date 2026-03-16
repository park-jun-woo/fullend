//ff:func feature=crosscheck type=util control=sequence
//ff:what SQL 줄이 제약조건(PRIMARY, UNIQUE 등)인지 확인
package crosscheck

import "strings"

func isConstraintLine(upper string) bool {
	return strings.HasPrefix(upper, "PRIMARY") || strings.HasPrefix(upper, "UNIQUE") ||
		strings.HasPrefix(upper, "CHECK") || strings.HasPrefix(upper, "CONSTRAINT") ||
		strings.HasPrefix(upper, "FOREIGN")
}
