//ff:func feature=manifest type=parser control=iteration dimension=1
//ff:what parseDDLContent — SQL 텍스트에서 CREATE TABLE 문을 라인 단위로 파싱
package ddl

import "strings"

func parseDDLContent(content string, tables map[string]*Table) {
	lines := strings.Split(content, "\n")
	var currentTable string
	for _, line := range lines {
		currentTable = parseDDLLine(line, currentTable, tables)
	}
	extractInserts(content, tables)
}
