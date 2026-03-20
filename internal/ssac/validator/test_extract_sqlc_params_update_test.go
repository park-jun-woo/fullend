//ff:func feature=symbol type=test control=iteration dimension=1 topic=sqlc
//ff:what UPDATE SQL에서 파라미터 추출 검증
package validator

import "testing"

func TestExtractSqlcParamsUpdate(t *testing.T) {
	sql := "UPDATE gigs SET status = $1 WHERE id = $2;"
	params := extractSqlcParams(sql)
	want := []string{"Status", "ID"}
	if len(params) != len(want) { t.Fatalf("expected %d params, got %d: %v", len(want), len(params), params) }
	for i, w := range want { if params[i] != w { t.Errorf("param[%d]: got %q, want %q", i, params[i], w) } }
}
