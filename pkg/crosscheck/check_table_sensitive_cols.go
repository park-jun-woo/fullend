//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkTableSensitiveCols — 테이블 컬럼 중 민감 패턴 매칭 검사
package crosscheck

func checkTableSensitiveCols(tableName string, columns map[string]string) []CrossError {
	var errs []CrossError
	for col := range columns {
		if matchesSensitivePattern(col) {
			errs = append(errs, CrossError{
				Rule: "X-61", Context: tableName + "." + col, Level: "WARNING",
				Message: "column matches sensitive pattern but has no @sensitive annotation",
			})
		}
	}
	return errs
}
