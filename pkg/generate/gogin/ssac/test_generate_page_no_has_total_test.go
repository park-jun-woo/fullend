//ff:func feature=ssac-gen type=test control=sequence
//ff:what Page[T] wrapper 사용 시 total 없이 단일 반환하는지 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func TestGeneratePageNoHasTotal(t *testing.T) {
	st := &rule.Ground{
		Models:     map[string]rule.ModelInfo{},
		Tables: map[string]rule.TableInfo{},
		Ops: map[string]rule.OperationInfo{
			"ListGigs": {Pagination: &rule.PaginationSpec{Style: "offset", DefaultLimit: 20, MaxLimit: 100}},
		},
	}
	sf := ssacparser.ServiceFunc{
		Name: "ListGigs", FileName: "list_gigs.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqGet, Model: "Gig.List", Inputs: map[string]string{"Query": "query"}, Result: &ssacparser.Result{Type: "Gig", Var: "gigPage", Wrapper: "Page"}},
			{Type: ssacparser.SeqResponse, Target: "gigPage"},
		},
	}
	code := mustGenerate(t, sf, st)
	// Page[T]이면 3-tuple 아니라 단일 반환
	assertNotContains(t, code, "total")
	assertContains(t, code, `gigPage, err :=`)
}
