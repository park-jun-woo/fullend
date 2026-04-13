//ff:func feature=ssac-gen type=test control=sequence
//ff:what x-sort 없을 때 SortConfig가 생성되지 않는지 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func TestGenerateSortNoXSort(t *testing.T) {
	st := &rule.Ground{
		Models:     map[string]rule.ModelInfo{},
		Tables: map[string]rule.TableInfo{},
		Ops: map[string]rule.OperationInfo{
			"ListGigs": {
				Pagination: &rule.PaginationSpec{Style: "offset", DefaultLimit: 20, MaxLimit: 100},
			},
		},
	}
	sf := ssacparser.ServiceFunc{
		Name: "ListGigs", FileName: "list_gigs.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqGet, Model: "Gig.List", Inputs: map[string]string{"Query": "query"}, Result: &ssacparser.Result{Type: "[]Gig", Var: "gigs"}},
			{Type: ssacparser.SeqResponse, Fields: map[string]string{"gigs": "gigs"}},
		},
	}
	code := mustGenerate(t, sf, st)
	// x-sort 없으면 SortConfig 없음
	assertNotContains(t, code, `Sort:`)
	assertNotContains(t, code, `SortConfig`)
}
