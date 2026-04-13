//ff:func feature=sqlc-parse type=util control=iteration dimension=1
//ff:what extractInsertParams — INSERT INTO table (col1, col2) VALUES ($1, $2) 에서 컬럼 순서 추출
package sqlc

import (
	"strings"

	"github.com/ettle/strcase"
)

func extractInsertParams(sql string) []string {
	parenStart := strings.IndexByte(sql, '(')
	if parenStart < 0 {
		return nil
	}
	parenEnd := strings.IndexByte(sql[parenStart:], ')')
	if parenEnd < 0 {
		return nil
	}
	colStr := sql[parenStart+1 : parenStart+parenEnd]
	var params []string
	for _, col := range strings.Split(colStr, ",") {
		col = strings.TrimSpace(col)
		if col != "" {
			params = append(params, strcase.ToGoPascal(col))
		}
	}
	return params
}
