//ff:func feature=crosscheck type=util control=sequence topic=ddl-coverage
//ff:what @archived 인라인 태그가 있는 컬럼을 ArchivedInfo에 기록
package crosscheck

import "strings"

func markArchivedColumn(trimmed, currentTable string, info *ArchivedInfo) {
	colParts := strings.Fields(trimmed)
	if len(colParts) < 2 {
		return
	}
	colName := colParts[0]
	if info.Columns[currentTable] == nil {
		info.Columns[currentTable] = make(map[string]bool)
	}
	info.Columns[currentTable][colName] = true
}
