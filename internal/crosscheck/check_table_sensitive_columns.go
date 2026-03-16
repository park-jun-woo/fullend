//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=sensitive
//ff:what 단일 테이블의 컬럼이 민감 패턴에 매치되지만 @sensitive 없는 경우 검증
package crosscheck

import "fmt"

func checkTableSensitiveColumns(tableName string, columnOrder []string, sensitiveCols, noSensitiveCols map[string]map[string]bool) []CrossError {
	var errs []CrossError
	for _, colName := range columnOrder {
		if isSensitiveAnnotated(tableName, colName, sensitiveCols) {
			continue
		}
		if isNoSensitiveAnnotated(tableName, colName, noSensitiveCols) {
			continue
		}
		if match := matchSensitivePattern(colName); match != "" {
			errs = append(errs, CrossError{
				Rule:       "DDL @sensitive",
				Context:    fmt.Sprintf("%s.%s", tableName, colName),
				Message:    fmt.Sprintf("column %q matches sensitive pattern %q but has no @sensitive annotation — will be exposed in JSON responses", colName, match),
				Level:      "WARNING",
				Suggestion: fmt.Sprintf("add -- @sensitive to the column definition in db/%s.sql to generate json:\"-\" tag", tableName),
			})
		}
	}
	return errs
}
