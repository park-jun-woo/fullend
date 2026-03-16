//ff:func feature=orchestrator type=rule control=sequence
//ff:what FK + DEFAULT 0 컬럼의 참조 테이블에 센티널 레코드가 있는지 검사한다

package orchestrator

import (
	"fmt"
	"regexp"
	"strings"
)

// checkSentinelRecord verifies that FK + DEFAULT 0 columns have a sentinel record
// in the referenced table. Returns an error message if missing, "" if OK.
func checkSentinelRecord(trimmed, upper, tableName, colName string, refRe *regexp.Regexp, tableContents map[string]string) string {
	if !strings.Contains(upper, "DEFAULT 0") || !strings.Contains(upper, "REFERENCES") {
		return ""
	}
	refMatch := refRe.FindStringSubmatch(trimmed)
	if refMatch == nil {
		return ""
	}
	refTable := refMatch[1]
	refContent, ok := tableContents[refTable]
	if !ok {
		return ""
	}
	if hasSentinelInsert(refContent, refTable) {
		return ""
	}
	return fmt.Sprintf("DDL: 테이블 %q 컬럼 %q — FK + DEFAULT 0이지만 참조 대상 %q에 id=0 센티널 레코드가 없습니다. INSERT INTO %s (id, ...) VALUES (0, ...) ON CONFLICT DO NOTHING; 을 추가하세요", tableName, colName, refTable, refTable)
}
