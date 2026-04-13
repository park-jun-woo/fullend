//ff:func feature=orchestrator type=rule control=sequence
//ff:what DDL 컬럼 라인에서 NOT NULL 누락 또는 FK DEFAULT 0 센티널 누락을 검사한다

package orchestrator

import (
	"fmt"
	"regexp"
	"strings"
)

// checkColumnLine inspects a single DDL line and returns an error message if
// the column violates NOT NULL or sentinel record rules. Returns "" if OK.
// skipSentinel: Phase018 auto_nobody_seed 활성 시 sentinel 검증 건너뜀.
func checkColumnLine(line, tableName string, colRe, refRe *regexp.Regexp, tableContents map[string]string, skipSentinel bool) string {
	trimmed := strings.TrimSpace(line)
	if isSkippableDDLLine(trimmed) {
		return ""
	}
	m := colRe.FindStringSubmatch(trimmed)
	if m == nil {
		return ""
	}
	colName := m[1]
	upper := strings.ToUpper(trimmed)
	if !strings.Contains(upper, "PRIMARY KEY") && !strings.Contains(upper, "NOT NULL") {
		return fmt.Sprintf("DDL: 테이블 %q 컬럼 %q — NOT NULL이 없습니다. NOT NULL DEFAULT 값을 지정하세요", tableName, colName)
	}
	if skipSentinel {
		return ""
	}
	if msg := checkSentinelRecord(trimmed, upper, tableName, colName, refRe, tableContents); msg != "" {
		return msg
	}
	return ""
}
