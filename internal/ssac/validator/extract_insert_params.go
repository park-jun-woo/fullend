//ff:func feature=symbol type=util
//ff:what INSERT INTO table (col1, col2) VALUES ($1, $2) 패턴에서 컬럼 순서를 추출한다
package validator

import (
	"strings"

	"github.com/ettle/strcase"
)

// extractInsertParams는 INSERT INTO table (col1, col2) VALUES ($1, $2) 패턴에서 컬럼 순서를 추출한다.
func extractInsertParams(sql string) []string {
	// 첫 번째 괄호 쌍 = 컬럼 목록
	parenStart := strings.IndexByte(sql, '(')
	if parenStart < 0 {
		return nil
	}
	parenEnd := strings.IndexByte(sql[parenStart:], ')')
	if parenEnd < 0 {
		return nil
	}
	colStr := sql[parenStart+1 : parenStart+parenEnd]
	cols := strings.Split(colStr, ",")

	var params []string
	for _, col := range cols {
		col = strings.TrimSpace(col)
		if col != "" {
			params = append(params, strcase.ToGoPascal(col))
		}
	}
	return params
}
