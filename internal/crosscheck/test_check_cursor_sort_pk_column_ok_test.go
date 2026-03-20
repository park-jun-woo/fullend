//ff:func feature=crosscheck type=rule control=sequence topic=openapi-ddl
//ff:what TestCheckCursorSort_PKColumn_OK: cursor 페이지네이션에 PK 컬럼으로 정렬하면 에러 없음 확인
package crosscheck

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func TestCheckCursorSort_PKColumn_OK(t *testing.T) {
	op := &openapi3.Operation{
		OperationID: "ListGigs",
		Extensions: map[string]any{
			"x-pagination": makeExt(map[string]any{"style": "cursor"}),
			"x-sort": makeExt(map[string]any{
				"allowed":   []string{"id"},
				"default":   "id",
				"direction": "desc",
			}),
		},
	}
	st := &ssacvalidator.SymbolTable{
		DDLTables: map[string]ssacvalidator.DDLTable{
			"gigs": {
				Columns:    map[string]string{"id": "int64", "title": "string"},
				PrimaryKey: []string{"id"},
			},
		},
	}
	errs := checkCursorSort(op, st, "GET /gigs (ListGigs)")
	if len(errs) != 0 {
		t.Errorf("expected no errors for PK column, got %v", errs)
	}
}
