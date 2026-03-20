//ff:func feature=symbol type=test control=iteration dimension=1 topic=ddl
//ff:what CREATE UNIQUE INDEX 파싱 검증
package validator

import "testing"

func TestDDLCreateUniqueIndex(t *testing.T) {
	tables := map[string]DDLTable{}
	ddl := "CREATE TABLE users (\n    id BIGSERIAL PRIMARY KEY,\n    email VARCHAR(255) NOT NULL\n);\nCREATE UNIQUE INDEX idx_users_email ON users (email);"
	parseDDLTables(ddl, tables)
	tbl := tables["users"]
	found := false
	for _, idx := range tbl.Indexes { if idx.IsUnique && idx.Name == "idx_users_email" { found = true } }
	if !found { t.Errorf("expected UNIQUE index idx_users_email, got indexes: %v", tbl.Indexes) }
}
