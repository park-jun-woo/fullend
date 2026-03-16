//ff:func feature=symbol type=util control=sequence topic=ddl
//ff:what DDLTable에 VARCHAR 길이를 설정
package validator

func applyVarcharLen(t *DDLTable, colName, colType string) {
	n := extractVarcharLen(colType)
	if n <= 0 {
		return
	}
	if t.VarcharLen == nil {
		t.VarcharLen = map[string]int{}
	}
	t.VarcharLen[colName] = n
}
