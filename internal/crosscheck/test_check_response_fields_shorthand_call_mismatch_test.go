//ff:func feature=crosscheck type=test control=iteration dimension=1 topic=ssac-openapi
//ff:what checkResponseFields: shorthand @response에서 funcspec JSON 필드명 불일치 검증

package crosscheck

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/funcspec"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestCheckResponseFields_ShorthandCallMismatch(t *testing.T) {
	// @response token + funcspec AccessToken(json:access_token) + OpenAPI AccessToken → ERROR
	doc := buildResponseDoc("Login", map[string]string{"AccessToken": "string"})
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
	// Should detect: "access_token" not in OpenAPI {AccessToken}
	foundErr := false
	for _, e := range errs {
		if e.Level != "WARNING" && contains(e.Message, "access_token") {
			foundErr = true
		}
	}
	if !foundErr {
		t.Errorf("expected ERROR for access_token mismatch, got: %+v", errs)
	}
}
