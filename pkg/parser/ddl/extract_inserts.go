//ff:func feature=manifest type=parser control=iteration dimension=2 topic=ddl
//ff:what extractInserts — 원본 SQL 내용에서 INSERT INTO ... VALUES 추출해 Table.Seeds 에 추가

package ddl

import (
	"regexp"
	"strings"
)

var reInsert = regexp.MustCompile(`(?is)INSERT\s+INTO\s+(\w+)\s*\(([^)]+)\)\s+VALUES\s*\(([^)]+)\)`)

// extractInserts scans raw SQL content for INSERT INTO <table> (cols) VALUES (vals)
// and appends each row to the target Table's Seeds field.
func extractInserts(content string, tables map[string]*Table) {
	matches := reInsert.FindAllStringSubmatch(content, -1)
	for _, m := range matches {
		cols := splitAndTrim(m[2], ",")
		vals := splitCSVLiterals(m[3])
		if len(cols) != len(vals) {
			continue
		}
		t := findTableCaseInsensitive(m[1], tables)
		if t == nil {
			continue
		}
		row := make(map[string]string, len(cols))
		for i, c := range cols {
			row[strings.ToLower(c)] = stripSQLQuotes(vals[i])
		}
		t.Seeds = append(t.Seeds, row)
	}
}
