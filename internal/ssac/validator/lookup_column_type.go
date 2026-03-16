//ff:func feature=ssac-validate type=util control=sequence topic=type-resolve
//ff:what DDL 심볼 테이블에서 테이블·컬럼의 Go 타입을 조회한다
package validator

func lookupColumnType(st *SymbolTable, tableName, colName string) string {
	table, ok := st.DDLTables[tableName]
	if !ok {
		return ""
	}
	goType, ok := table.Columns[colName]
	if !ok {
		return ""
	}
	return goType
}
