//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=ddl-coverage
//ff:what SQL 텍스트에서 @archived 테이블·컬럼 태그를 추출
package crosscheck

import "strings"

func parseArchivedSQL(content string, info *ArchivedInfo) {
	lines := strings.Split(content, "\n")
	prevLineArchived := false
	var currentTable string

	for _, line := range lines {
		prevLineArchived, currentTable = parseArchivedLine(line, prevLineArchived, currentTable, info)
	}
}
