//ff:func feature=symbol type=util
//ff:what CREATE INDEX 문을 파싱한다
package validator

import "strings"

// parseCreateIndex는 CREATE INDEX 문을 파싱한다.
// e.g. "CREATE INDEX idx_name ON tablename (col1, col2);"
func parseCreateIndex(line string, tables map[string]DDLTable) {
	upper := strings.ToUpper(line)
	onIdx := strings.Index(upper, " ON ")
	if onIdx < 0 {
		return
	}

	// 인덱스 이름: CREATE [UNIQUE] INDEX idx_name ON ...
	parts := strings.Fields(line[:onIdx])
	idxName := ""
	for i, p := range parts {
		if strings.ToUpper(p) == "INDEX" && i+1 < len(parts) {
			idxName = parts[i+1]
			break
		}
	}

	// ON tablename (col1, col2)
	after := strings.TrimSpace(line[onIdx+4:])
	afterParts := strings.SplitN(after, "(", 2)
	if len(afterParts) < 2 {
		return
	}

	tableName := strings.TrimSpace(afterParts[0])
	colsPart := strings.TrimSuffix(strings.TrimSpace(afterParts[1]), ");")
	colsPart = strings.TrimSuffix(colsPart, ")")

	var cols []string
	for _, c := range strings.Split(colsPart, ",") {
		c = strings.TrimSpace(c)
		if c != "" {
			cols = append(cols, c)
		}
	}

	isUnique := strings.Contains(strings.ToUpper(line), "UNIQUE")
	if t, ok := tables[tableName]; ok && len(cols) > 0 {
		t.Indexes = append(t.Indexes, Index{Name: idxName, Columns: cols, IsUnique: isUnique})
		tables[tableName] = t
	}
}
