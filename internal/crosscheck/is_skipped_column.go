//ff:func feature=crosscheck type=util control=sequence
//ff:what 컬럼이 민감 정보 또는 x-include로 스킵 대상인지 확인
package crosscheck

// isSkippedColumn checks if a column should be skipped in missing property checks.
func isSkippedColumn(tableName, colName string, sensitiveCols map[string]map[string]bool, xIncludeFields map[string]bool) bool {
	if sensitiveCols != nil {
		if cols, ok := sensitiveCols[tableName]; ok && cols[colName] {
			return true
		}
	}
	if matchesSensitivePattern(colName) {
		return true
	}
	if xIncludeFields[colName] {
		return true
	}
	return false
}
