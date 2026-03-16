//ff:func feature=symbol type=util control=iteration dimension=1
//ff:what CREATE TABLE 문에서 테이블명을 추출하고 tables에 등록한다
package validator

import "strings"

func extractAndRegisterTable(line string, tables map[string]DDLTable) string {
	parts := strings.Fields(line)
	for i, p := range parts {
		pu := strings.ToUpper(p)
		if pu != "TABLE" || i+1 >= len(parts) {
			continue
		}
		name := strings.Trim(parts[i+1], "( ")
		if name != "" {
			tables[name] = DDLTable{Columns: make(map[string]string)}
		}
		return name
	}
	return ""
}
