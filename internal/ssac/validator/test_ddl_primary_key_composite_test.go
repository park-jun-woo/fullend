//ff:func feature=symbol type=test control=sequence topic=ddl
//ff:what DDL 복합 PRIMARY KEY 파싱 검증
package validator

import "testing"

func TestDDLPrimaryKeyComposite(t *testing.T) {
	tables := map[string]DDLTable{}
	ddl := "CREATE TABLE enrollments (\n    user_id BIGINT NOT NULL,\n    course_id BIGINT NOT NULL,\n    PRIMARY KEY (user_id, course_id)\n);"
	parseDDLTables(ddl, tables)
	tbl := tables["enrollments"]
	if len(tbl.PrimaryKey) != 2 { t.Fatalf("expected 2 PK columns, got %d: %v", len(tbl.PrimaryKey), tbl.PrimaryKey) }
	if tbl.PrimaryKey[0] != "user_id" || tbl.PrimaryKey[1] != "course_id" { t.Errorf("expected PrimaryKey=[user_id, course_id], got %v", tbl.PrimaryKey) }
}
