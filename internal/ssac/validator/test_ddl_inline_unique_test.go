//ff:func feature=symbol type=test control=iteration dimension=1 topic=ddl
//ff:what DDL 인라인 UNIQUE 파싱 검증
package validator

import "testing"

func TestDDLInlineUnique(t *testing.T) {
	tables := map[string]DDLTable{}
	ddl := "CREATE TABLE users (\n    id BIGSERIAL PRIMARY KEY,\n    email VARCHAR(255) NOT NULL UNIQUE\n);"
	parseDDLTables(ddl, tables)
	tbl := tables["users"]
	found := false
	for _, idx := range tbl.Indexes { if idx.IsUnique && len(idx.Columns) == 1 && idx.Columns[0] == "email" { found = true } }
	if !found { t.Errorf("expected UNIQUE index for email, got indexes: %v", tbl.Indexes) }
}
