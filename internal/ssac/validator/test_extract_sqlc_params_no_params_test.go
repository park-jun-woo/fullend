//ff:func feature=symbol type=test control=sequence topic=sqlc
//ff:what 파라미터 없는 SQL에서 빈 결과 검증
package validator

import "testing"

func TestExtractSqlcParamsNoParams(t *testing.T) {
	sql := "SELECT * FROM gigs ORDER BY created_at DESC;"
	params := extractSqlcParams(sql)
	if len(params) != 0 { t.Errorf("expected no params, got %v", params) }
}
