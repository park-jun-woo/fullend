//ff:func feature=crosscheck type=rule control=sequence topic=openapi-ddl
//ff:what TestCheckCursorSort_UniqueDefault_OK: cursor 페이지네이션에 UNIQUE 정렬 컬럼이면 에러 없음 확인
package crosscheck

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func TestCheckCursorSort_UniqueDefault_OK(t *testing.T) {
	op := &openapi3.Operation{
		OperationID: "ListOrders",
		Extensions: map[string]any{
			"x-pagination": makeExt(map[string]any{"style": "cursor"}),
			"x-sort": makeExt(map[string]any{
				"allowed":   []string{"order_id"},
				"default":   "order_id",
				"direction": "desc",
			}),
		},
	}
	st := &ssacvalidator.SymbolTable{
		DDLTables: map[string]ssacvalidator.DDLTable{
			"orders": {
				Columns:    map[string]string{"id": "int64", "order_id": "string"},
				PrimaryKey: []string{"id"},
				Indexes: []ssacvalidator.Index{
					{Name: "uniq_order_id", Columns: []string{"order_id"}, IsUnique: true},
				},
			},
		},
	}
	errs := checkCursorSort(op, st, "GET /orders (ListOrders)")
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}
