//ff:func feature=rule type=loader control=sequence
//ff:what populateTables 검증 — DDL 테이블 → Ground.Tables 복사
package ground

import (
	"testing"

	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/parser/ddl"
)

func TestPopulateTables(t *testing.T) {
	g := newGround()
	fs := &fullend.Fullstack{
		DDLTables: []ddl.Table{
			{
				Name:        "users",
				Columns:     map[string]string{"id": "int64", "email": "string"},
				ColumnOrder: []string{"id", "email"},
			},
			{
				Name:        "courses",
				Columns:     map[string]string{"id": "int64", "title": "string"},
				ColumnOrder: []string{"id", "title"},
			},
		},
	}
	populateTables(g, fs)

	users, ok := g.Tables["users"]
	if !ok {
		t.Fatal("users not populated")
	}
	if users.Columns["email"] != "string" {
		t.Errorf("email type: got %q", users.Columns["email"])
	}
	if len(users.ColumnOrder) != 2 || users.ColumnOrder[0] != "id" || users.ColumnOrder[1] != "email" {
		t.Errorf("column order: got %v", users.ColumnOrder)
	}
	if _, ok := g.Tables["courses"]; !ok {
		t.Error("courses not populated")
	}
}
