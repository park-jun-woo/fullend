//ff:func feature=symbol type=test control=iteration dimension=1 topic=ddl
//ff:what DDL UNIQUE 제약조건 파싱 검증
package validator

import "testing"

func TestDDLUniqueConstraint(t *testing.T) {
	tables := map[string]DDLTable{}
	ddl := "CREATE TABLE reservations (\n    id BIGSERIAL PRIMARY KEY,\n    room_id BIGINT NOT NULL,\n    start_time TIMESTAMP NOT NULL,\n    UNIQUE (room_id, start_time)\n);"
	parseDDLTables(ddl, tables)
	tbl := tables["reservations"]
	found := false
	for _, idx := range tbl.Indexes { if idx.IsUnique && len(idx.Columns) == 2 && idx.Columns[0] == "room_id" && idx.Columns[1] == "start_time" { found = true } }
	if !found { t.Errorf("expected UNIQUE index for (room_id, start_time), got indexes: %v", tbl.Indexes) }
}
