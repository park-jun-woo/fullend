//ff:func feature=sqlc-parse type=util control=sequence
//ff:what extractParams — SQL 본문에서 $N ↔ 컬럼명을 추출해 순서대로 PascalCase 반환
package sqlc

import "strings"

// extractParams는 SQL 본문에서 $N ↔ 컬럼명 매핑을 추출해 $1, $2, ... 순서의 PascalCase 파라미터명을 반환한다.
func extractParams(sql string) []string {
	if !strings.Contains(sql, "$") {
		return nil
	}
	if strings.Contains(strings.ToUpper(sql), "INSERT") {
		if params := extractInsertParams(sql); len(params) > 0 {
			return params
		}
	}
	return extractWhereSetParams(sql)
}
