//ff:func feature=crosscheck type=rule control=sequence topic=openapi-ddl
//ff:what TestCheckCursorSort_OffsetIgnored: offset 페이지네이션은 cursor 정렬 검증 대상이 아님 확인
package crosscheck

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func TestCheckCursorSort_OffsetIgnored(t *testing.T) {
	op := &openapi3.Operation{
		OperationID: "ListGigs",
		Extensions: map[string]any{
			"x-pagination": makeExt(map[string]any{"style": "offset"}),
			"x-sort": makeExt(map[string]any{
				"allowed":   []string{"created_at", "budget"},
				"default":   "created_at",
				"direction": "desc",
			}),
		},
	}
	st := &ssacvalidator.SymbolTable{
		DDLTables: map[string]ssacvalidator.DDLTable{
			"gigs": {Columns: map[string]string{"id": "int64"}},
		},
	}
	errs := checkCursorSort(op, st, "GET /gigs (ListGigs)")
	if len(errs) != 0 {
		t.Errorf("expected no errors for offset mode, got %v", errs)
	}
}
