//ff:func feature=ssac-gen type=test control=sequence
//ff:what @response target 직접 반환 시 gin.H 래핑 없이 직접 반환하는지 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestGenerateResponseDirect(t *testing.T) {
	sf := parser.ServiceFunc{
		Name: "ListGigs", FileName: "list_gigs.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Gig.List", Inputs: map[string]string{"Query": "query"}, Result: &parser.Result{Type: "Gig", Var: "gigPage", Wrapper: "Page"}},
			{Type: parser.SeqResponse, Target: "gigPage"},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, `c.JSON(__RESPONSE_STATUS__, gigPage)`)
	assertNotContains(t, code, `c.JSON(__RESPONSE_STATUS__, gin.H`)
	assertNotContains(t, code, `pagination`)
}
