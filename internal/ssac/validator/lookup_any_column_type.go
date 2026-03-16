//ff:func feature=symbol type=util control=iteration dimension=1 topic=ddl
//ff:what 모든 DDL 테이블에서 해당 컬럼명의 Go 타입을 찾는다
package validator

// lookupAnyColumnType는 모든 DDL 테이블에서 해당 컬럼명의 Go 타입을 찾는다.
func lookupAnyColumnType(st *SymbolTable, colName string) string {
	for _, table := range st.DDLTables {
		if goType, ok := table.Columns[colName]; ok {
			return goType
		}
	}
	return ""
}
