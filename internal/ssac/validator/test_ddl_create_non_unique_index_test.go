//ff:func feature=symbol type=test control=iteration dimension=1 topic=ddl
//ff:what CREATE INDEX(비-UNIQUE) 파싱 검증
package validator

import "testing"

func TestDDLCreateNonUniqueIndex(t *testing.T) {
	tables := map[string]DDLTable{}
	ddl := "CREATE TABLE orders (\n    id BIGSERIAL PRIMARY KEY,\n    status VARCHAR(50) NOT NULL\n);\nCREATE INDEX idx_orders_status ON orders (status);"
	parseDDLTables(ddl, tables)
	tbl := tables["orders"]
	for _, idx := range tbl.Indexes { if idx.Name == "idx_orders_status" && idx.IsUnique { t.Errorf("expected non-unique index, got IsUnique=true") } }
}
