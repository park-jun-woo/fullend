//ff:func feature=crosscheck type=test control=iteration dimension=1 topic=ssac-openapi
//ff:what checkResponseFields: shorthand @response에서 DDL 컬럼과 OpenAPI 필드 일치 검증

package crosscheck

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func TestCheckResponseFields_ShorthandDDLMatch(t *testing.T) {
	// @response user + DDL columns [id, email, name] + OpenAPI [id, email, name] → pass
	doc := buildResponseDoc("GetUser", map[string]string{"id": "integer", "email": "string", "name": "string"})
	st := &ssacvalidator.SymbolTable{
		DDLTables: map[string]ssacvalidator.DDLTable{
			"users": {Columns: map[string]string{"id": "int64", "email": "string", "name": "string"}},
		},
	}
	funcs := []ssacparser.ServiceFunc{{
		Name:     "GetUser",
		FileName: "get_user.ssac",
		Sequences: []ssacparser.Sequence{
			{Type: "get", Result: &ssacparser.Result{Type: "User", Var: "user"}},
			{Type: "response", Target: "user"},
		},
	}}

	errs := checkResponseFields(funcs, st, doc, nil)
	for _, e := range errs {
		if e.Level != "WARNING" {
			t.Errorf("unexpected ERROR: %+v", e)
		}
	}
}
