//ff:func feature=symbol type=util control=sequence topic=ddl
//ff:what CONSTRAINT FK 라인을 파싱하여 DDLTable에 추가한다
package validator

func appendConstraintFK(line, currentTable string, tables map[string]DDLTable) {
	fk, ok := parseConstraintFK(line)
	if !ok {
		return
	}
	t, exists := tables[currentTable]
	if !exists {
		return
	}
	t.ForeignKeys = append(t.ForeignKeys, fk)
	tables[currentTable] = t
}
