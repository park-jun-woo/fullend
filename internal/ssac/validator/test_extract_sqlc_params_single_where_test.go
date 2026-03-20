//ff:func feature=symbol type=test control=sequence topic=sqlc
//ff:what 단일 WHERE 파라미터 추출 검증
package validator

import "testing"

func TestExtractSqlcParamsSingleWhere(t *testing.T) {
	sql := "SELECT * FROM users WHERE id = $1;"
	params := extractSqlcParams(sql)
	want := []string{"ID"}
	if len(params) != len(want) { t.Fatalf("expected %d params, got %d: %v", len(want), len(params), params) }
	if params[0] != want[0] { t.Errorf("param[0]: got %q, want %q", params[0], want[0]) }
}
