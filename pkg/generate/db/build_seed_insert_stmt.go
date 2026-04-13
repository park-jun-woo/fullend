//ff:func feature=gen-gogin type=generator control=iteration dimension=1 topic=ddl
//ff:what buildSeedInsertStmt — 테이블 NOT NULL 컬럼들에 대한 placeholder INSERT 생성

package db

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/fullend/pkg/parser/ddl"
)

// buildSeedInsertStmt emits: INSERT INTO <t> (cols) VALUES (vals) ON CONFLICT DO NOTHING;
// DEFAULT 가 있는 컬럼(id 제외) 은 skip — DB 가 채움.
func buildSeedInsertStmt(t *ddl.Table, idVal int64) string {
	var cols, vals []string
	for _, col := range t.ColumnOrder {
		if _, hasDefault := t.Defaults[col]; hasDefault && col != "id" {
			continue
		}
		cols = append(cols, col)
		vals = append(vals, seedValueFor(t, col, idVal))
	}
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) ON CONFLICT DO NOTHING;",
		t.Name, strings.Join(cols, ", "), strings.Join(vals, ", "))
}
