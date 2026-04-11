//ff:func feature=manifest type=parser control=sequence
//ff:what parseDDLLine — DDL 한 줄을 파싱하여 테이블/컬럼/제약조건 반영
package ddl

import "strings"

func parseDDLLine(line string, currentTable string, tables map[string]*Table) string {
	line = strings.TrimSpace(line)
	upper := strings.ToUpper(line)

	if strings.HasPrefix(upper, "CREATE INDEX") || strings.HasPrefix(upper, "CREATE UNIQUE INDEX") {
		parseCreateIndex(line, tables)
		return currentTable
	}
	if strings.HasPrefix(upper, "CREATE TABLE") {
		return extractTableName(line, tables)
	}
	if currentTable == "" || strings.HasPrefix(line, ")") {
		if strings.HasPrefix(line, ")") {
			return ""
		}
		return currentTable
	}
	t := tables[currentTable]
	if t == nil {
		return currentTable
	}
	dispatchConstraint(line, upper, t, tables)
	return currentTable
}
