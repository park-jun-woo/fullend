//ff:func feature=symbol type=test control=iteration dimension=1 topic=sqlc
//ff:what INSERT SQL에서 파라미터 추출 검증
package validator

import "testing"

func TestExtractSqlcParamsInsert(t *testing.T) {
	sql := "INSERT INTO users (email, password_hash, name) VALUES ($1, $2, $3) RETURNING *;"
	params := extractSqlcParams(sql)
	want := []string{"Email", "PasswordHash", "Name"}
	if len(params) != len(want) { t.Fatalf("expected %d params, got %d: %v", len(want), len(params), params) }
	for i, w := range want { if params[i] != w { t.Errorf("param[%d]: got %q, want %q", i, params[i], w) } }
}
