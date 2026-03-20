//ff:func feature=crosscheck type=test control=iteration dimension=1 topic=ssac-openapi
//ff:what checkResponseFields: shorthand @response에서 funcspec JSON 필드명 일치 검증

package crosscheck

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/funcspec"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestCheckResponseFields_ShorthandCallMatch(t *testing.T) {
	// @response token + funcspec AccessToken(json:access_token) + OpenAPI access_token → pass
	doc := buildResponseDoc("Login", map[string]string{"access_token": "string"})
	funcSpecs := []funcspec.FuncSpec{{
		Package: "auth",
		Name:    "issueToken",
		ResponseFields: []funcspec.Field{
			{Name: "AccessToken", Type: "string", JSONName: "access_token"},
		},
	}}
	funcs := []ssacparser.ServiceFunc{{
		Name:     "Login",
		FileName: "login.ssac",
		Sequences: []ssacparser.Sequence{
			{Type: "call", Result: &ssacparser.Result{Type: "auth.IssueTokenResponse", Var: "token"}},
			{Type: "response", Target: "token"},
		},
	}}

	errs := checkResponseFields(funcs, nil, doc, funcSpecs)
	for _, e := range errs {
		if e.Level != "WARNING" {
			t.Errorf("unexpected ERROR: %+v", e)
		}
	}
}
