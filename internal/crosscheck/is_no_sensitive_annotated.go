//ff:func feature=crosscheck type=util control=sequence topic=sensitive
//ff:what 컬럼이 @nosensitive로 표시되어 있는지 확인
package crosscheck

func isNoSensitiveAnnotated(tableName, colName string, noSensitiveCols map[string]map[string]bool) bool {
	if noSensitiveCols == nil {
		return false
	}
	cols, ok := noSensitiveCols[tableName]
	return ok && cols[colName]
}
