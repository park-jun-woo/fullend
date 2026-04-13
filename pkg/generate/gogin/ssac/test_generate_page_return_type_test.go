//ff:func feature=ssac-gen type=test control=sequence
//ff:what Page[T] wrapper 사용 시 인터페이스 반환 타입을 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func TestGeneratePageReturnType(t *testing.T) {
	st := &rule.Ground{
		Models: map[string]rule.ModelInfo{
			"Gig": {Methods: map[string]rule.MethodInfo{
				"List": {Cardinality: "many"},
			}},
		},
		Tables: map[string]rule.TableInfo{},
		Ops: map[string]rule.OperationInfo{},
	}
	funcs := []ssacparser.ServiceFunc{{
		Name: "ListGigs", FileName: "list_gigs.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqGet, Model: "Gig.List", Inputs: map[string]string{"Query": "query"}, Result: &ssacparser.Result{Type: "Gig", Var: "gigPage", Wrapper: "Page"}},
			{Type: ssacparser.SeqResponse, Target: "gigPage"},
		},
	}}

	outDir := t.TempDir()
	if err := GenerateModelInterfaces(funcs, st, outDir); err != nil {
		t.Fatal(err)
	}

	data, err := readFile(t, outDir+"/model/models_gen.go")
	if err != nil {
		t.Fatal(err)
	}
	assertContains(t, data, "(*pagination.Page[Gig], error)")
	assertContains(t, data, "pagination")
}
