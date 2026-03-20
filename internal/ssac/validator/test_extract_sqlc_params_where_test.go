//ff:func feature=symbol type=test control=iteration dimension=1 topic=sqlc
//ff:what WHERE절 SQL에서 파라미터 추출 검증
package validator

import "testing"

func TestExtractSqlcParamsWhere(t *testing.T) {
	sql := "SELECT * FROM reservations WHERE room_id = $1 AND end_at > $2 AND start_at < $3 LIMIT 1;"
	params := extractSqlcParams(sql)
	want := []string{"RoomID", "EndAt", "StartAt"}
	if len(params) != len(want) { t.Fatalf("expected %d params, got %d: %v", len(want), len(params), params) }
	for i, w := range want { if params[i] != w { t.Errorf("param[%d]: got %q, want %q", i, params[i], w) } }
}
