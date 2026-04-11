//ff:func feature=manifest type=util control=sequence
//ff:what extractTableName — CREATE TABLE 문에서 테이블명을 추출하고 등록
package ddl

import "strings"

func extractTableName(line string, tables map[string]*Table) string {
	parts := strings.Fields(line)
	idx := findTableKeyword(parts)
	if idx < 0 || idx+1 >= len(parts) {
		return ""
	}
	name := strings.Trim(parts[idx+1], "( ")
	if name != "" {
		tables[name] = &Table{Name: name, Columns: make(map[string]string)}
	}
	return name
}
