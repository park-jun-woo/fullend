//ff:func feature=crosscheck type=rule control=sequence topic=openapi-ddl
//ff:what TestCheckCursorSort_MultipleAllowed_ERROR: cursor 페이지네이션에 다중 정렬 컬럼이면 ERROR 확인
package crosscheck

import (
	"strings"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func TestCheckCursorSort_MultipleAllowed_ERROR(t *testing.T) {
	op := &openapi3.Operation{
		OperationID: "ListGigs",
		Extensions: map[string]any{
			"x-pagination": makeExt(map[string]any{"style": "cursor"}),
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
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(errs), errs)
	}
	if !strings.Contains(errs[0].Message, "런타임 정렬 전환") {
		t.Errorf("expected runtime sort error, got %q", errs[0].Message)
	}
}
