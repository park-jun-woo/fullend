//ff:func feature=crosscheck type=rule control=sequence topic=openapi-ddl
//ff:what TestCheckCursorSort_NoSort_OK: cursor 페이지네이션에 x-sort 없으면 에러 없음 확인
package crosscheck

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func TestCheckCursorSort_NoSort_OK(t *testing.T) {
	op := &openapi3.Operation{
		OperationID: "ListGigs",
		Extensions: map[string]any{
			"x-pagination": makeExt(map[string]any{"style": "cursor", "defaultLimit": 20}),
		},
	}
	st := &ssacvalidator.SymbolTable{
		DDLTables: map[string]ssacvalidator.DDLTable{
			"gigs": {Columns: map[string]string{"id": "int64"}},
		},
	}
	errs := checkCursorSort(op, st, "GET /gigs (ListGigs)")
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}
