//ff:func feature=symbol type=util control=sequence
//ff:what SQL 본문에서 $N ↔ 컬럼명 매핑을 추출하여 파라미터명을 반환한다
package validator

import "strings"

// extractSqlcParams는 SQL 본문에서 $N ↔ 컬럼명 매핑을 추출하여 $1, $2, ... 순서의 PascalCase 파라미터명을 반환한다.
func extractSqlcParams(sql string) []string {
	if !strings.Contains(sql, "$") {
		return nil
	}

	upper := strings.ToUpper(sql)
	if strings.Contains(upper, "INSERT") {
		if params := extractInsertParams(sql); len(params) > 0 {
			return params
		}
	}
	return extractWhereSetParams(sql)
}
