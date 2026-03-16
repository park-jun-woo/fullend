//ff:func feature=crosscheck type=util control=sequence
//ff:what SQL 한 줄에서 @sensitive/@nosensitive 태그를 감지하고 테이블·컬럼 정보를 갱신
package crosscheck

import "strings"

func parseSensitiveLine(line, currentTable string, sensitive, nosensitive map[string]map[string]bool) string {
	trimmed := strings.TrimSpace(line)
	upper := strings.ToUpper(trimmed)

	if strings.HasPrefix(upper, "CREATE TABLE") {
		return extractTableName(trimmed)
	}

	if strings.HasPrefix(trimmed, ")") {
		return ""
	}

	if currentTable == "" {
		return currentTable
	}

	colParts := strings.Fields(trimmed)
	if len(colParts) < 2 {
		return currentTable
	}
	colName := colParts[0]
	upperFirst := strings.ToUpper(colName)
	if upperFirst == "PRIMARY" || upperFirst == "UNIQUE" || upperFirst == "CHECK" ||
		upperFirst == "CONSTRAINT" || upperFirst == "FOREIGN" || upperFirst == "--" {
		return currentTable
	}

	if strings.Contains(line, "@nosensitive") {
		if nosensitive[currentTable] == nil {
			nosensitive[currentTable] = make(map[string]bool)
		}
		nosensitive[currentTable][colName] = true
	} else if strings.Contains(line, "@sensitive") {
		if sensitive[currentTable] == nil {
			sensitive[currentTable] = make(map[string]bool)
		}
		sensitive[currentTable][colName] = true
	}

	return currentTable
}
