//ff:func feature=crosscheck type=util control=sequence
//ff:what SQL 한 줄에서 @archived 태그를 감지하고 테이블·컬럼 정보를 갱신
package crosscheck

import "strings"

func parseArchivedLine(line string, prevArchived bool, currentTable string, info *ArchivedInfo) (bool, string) {
	trimmed := strings.TrimSpace(line)
	upper := strings.ToUpper(trimmed)

	if strings.HasPrefix(trimmed, "--") && strings.Contains(trimmed, "@archived") {
		return true, currentTable
	}

	if strings.HasPrefix(upper, "CREATE TABLE") {
		tableName := extractTableName(trimmed)
		if prevArchived && tableName != "" {
			info.Tables[tableName] = true
		}
		return false, tableName
	}

	if currentTable == "" {
		return false, currentTable
	}

	if strings.HasPrefix(trimmed, ")") {
		return false, ""
	}

	if isConstraintLine(upper) || trimmed == "" {
		return false, currentTable
	}

	if strings.Contains(line, "-- @archived") {
		markArchivedColumn(trimmed, currentTable, info)
	}

	return false, currentTable
}
