//ff:func feature=symbol type=parser control=iteration dimension=1 topic=ddl
//ff:what CREATE TABLE 문에서 컬럼명, 타입, FK, 인덱스를 추출한다
package validator

import "strings"

// parseDDLTables는 CREATE TABLE 문에서 컬럼명, 타입, FK, 인덱스를 추출한다.
func parseDDLTables(content string, tables map[string]DDLTable) {
	lines := strings.Split(content, "\n")
	var currentTable string

	for _, line := range lines {
		currentTable = parseDDLLine(line, currentTable, tables)
	}
}
