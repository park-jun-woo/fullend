//ff:func feature=symbol type=test control=sequence topic=ddl
//ff:what DDL PRIMARY KEY 파싱 검증
package validator

import "testing"

func TestDDLPrimaryKey(t *testing.T) {
	tables := map[string]DDLTable{}
	ddl := "CREATE TABLE users (\n    id BIGSERIAL PRIMARY KEY,\n    email VARCHAR(255) NOT NULL\n);"
	parseDDLTables(ddl, tables)
	tbl := tables["users"]
	if len(tbl.PrimaryKey) != 1 || tbl.PrimaryKey[0] != "id" { t.Errorf("expected PrimaryKey=[id], got %v", tbl.PrimaryKey) }
}
