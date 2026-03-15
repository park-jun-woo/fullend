//ff:func feature=orchestrator type=util
//ff:what 센티널 레코드(id=0) INSERT 존재 여부 확인
package orchestrator

import (
	"regexp"
	"strings"
)

// hasSentinelInsert checks if the DDL content contains an INSERT with id=0 for the given table.
func hasSentinelInsert(content, tableName string) bool {
	upper := strings.ToUpper(content)
	// INSERT INTO <table> ... VALUES (0, ...)
	insertRe := regexp.MustCompile(`(?i)INSERT\s+INTO\s+` + tableName + `\b`)
	if !insertRe.MatchString(content) {
		return false
	}
	// VALUES 절에서 첫 번째 값이 0인지 확인.
	valuesRe := regexp.MustCompile(`(?i)VALUES\s*\(\s*0\s*,`)
	idx := insertRe.FindStringIndex(upper)
	if idx == nil {
		return false
	}
	return valuesRe.MatchString(content[idx[0]:])
}
