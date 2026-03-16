//ff:func feature=symbol type=util control=sequence topic=ddl
//ff:what UNIQUE 제약 라인에서 unique index를 추가한다
package validator

import "strings"

func appendUniqueIndex(line, currentTable string, tables map[string]DDLTable) {
	t, ok := tables[currentTable]
	if !ok {
		return
	}
	cols := extractParenColumns(line)
	if len(cols) == 0 {
		return
	}
	t.Indexes = append(t.Indexes, Index{Name: "unique_" + strings.Join(cols, "_"), Columns: cols, IsUnique: true})
	tables[currentTable] = t
}
