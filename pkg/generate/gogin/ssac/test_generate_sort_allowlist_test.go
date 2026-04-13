//ff:func feature=ssac-gen type=test control=sequence
//ff:what x-sort allowlist를 SortConfig로 생성하는지 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func TestGenerateSortAllowlist(t *testing.T) {
	st := &validator.SymbolTable{
		Models:     map[string]validator.ModelSymbol{},
		DDLTables:  map[string]validator.DDLTable{},
		Operations: map[string]validator.OperationSymbol{
			"ListGigs": {
				XPagination: &validator.XPagination{Style: "offset", DefaultLimit: 20, MaxLimit: 100},
				XSort:       &validator.XSort{Allowed: []string{"created_at", "title", "price"}, Default: "created_at"},
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
	// model.ParseQueryOpts로 sort allowlist 포함 config 생성
	assertContains(t, code, `model.ParseQueryOpts(c, model.QueryOptsConfig{`)
	assertContains(t, code, `&model.SortConfig{`)
	assertContains(t, code, `"created_at"`)
	assertContains(t, code, `"title"`)
	assertContains(t, code, `"price"`)
	assertContains(t, code, `Default: "created_at"`)
	// 수동 파싱 패턴 없어야 함
	assertNotContains(t, code, `allowedSort`)
	assertNotContains(t, code, `c.Query("sort")`)
}
