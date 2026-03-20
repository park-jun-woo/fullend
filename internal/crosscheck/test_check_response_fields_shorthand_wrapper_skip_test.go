//ff:func feature=crosscheck type=test control=sequence topic=ssac-openapi
//ff:what checkResponseFields: Page 래퍼 타입 @response는 필드 검사 스킵 검증

package crosscheck

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func TestCheckResponseFields_ShorthandWrapperSkip(t *testing.T) {
	// @response gigPage with Page wrapper → should be skipped (no errors)
	doc := buildResponseDoc("ListGigs", map[string]string{"items": "array", "total": "integer"})
	st := &ssacvalidator.SymbolTable{
		DDLTables: map[string]ssacvalidator.DDLTable{
			"gigs": {Columns: map[string]string{"id": "int64", "title": "string"}},
		},
	}
	funcs := []ssacparser.ServiceFunc{{
		Name:     "ListGigs",
		FileName: "list_gigs.ssac",
		Sequences: []ssacparser.Sequence{
			{Type: "get", Result: &ssacparser.Result{Type: "Gig", Var: "gigPage", Wrapper: "Page"}},
			{Type: "response", Target: "gigPage"},
		},
	}}

	errs := checkResponseFields(funcs, st, doc, nil)
	if len(errs) != 0 {
		t.Errorf("expected 0 errors for wrapper type, got %d: %+v", len(errs), errs)
	}
}
