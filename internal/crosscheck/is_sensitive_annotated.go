//ff:func feature=crosscheck type=util control=sequence topic=sensitive
//ff:what 컬럼이 @sensitive로 표시되어 있는지 확인
package crosscheck

func isSensitiveAnnotated(tableName, colName string, sensitiveCols map[string]map[string]bool) bool {
	if sensitiveCols == nil {
		return false
	}
	cols, ok := sensitiveCols[tableName]
	return ok && cols[colName]
}
