//ff:func feature=symbol type=util control=sequence
//ff:what PRIMARY KEY 라인에서 PK 컬럼을 추출하여 설정한다
package validator

func setPrimaryKey(line, currentTable string, tables map[string]DDLTable) {
	t, ok := tables[currentTable]
	if !ok {
		return
	}
	t.PrimaryKey = extractParenColumns(line)
	tables[currentTable] = t
}
