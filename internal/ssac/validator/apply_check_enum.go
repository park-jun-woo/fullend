//ff:func feature=symbol type=util control=sequence topic=ddl
//ff:what CHECK enum 값을 DDLTable에 적용 (독립/인라인 CHECK 공용)
package validator

// applyCheckEnum parses a CHECK line and applies enum values to the table.
// If colName is empty, it uses the column name from the CHECK clause itself (standalone CHECK).
// If colName is non-empty, it uses the provided column name (inline CHECK in column line).
func applyCheckEnum(line, colName, currentTable string, tables map[string]DDLTable) {
	col, vals := parseCheckEnum(line)
	if colName != "" {
		col = colName
	}
	if col == "" || len(vals) == 0 {
		return
	}
	t, ok := tables[currentTable]
	if !ok {
		return
	}
	if t.CheckEnums == nil {
		t.CheckEnums = map[string][]string{}
	}
	t.CheckEnums[col] = vals
	tables[currentTable] = t
}
