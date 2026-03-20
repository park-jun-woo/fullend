//ff:func feature=ssac-gen type=test control=sequence
//ff:what Page[T] wrapper 사용 시 인터페이스 반환 타입을 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func TestGeneratePageReturnType(t *testing.T) {
	st := &validator.SymbolTable{
		Models: map[string]validator.ModelSymbol{
			"Gig": {Methods: map[string]validator.MethodInfo{
				"List": {Cardinality: "many"},
			}},
		},
		DDLTables:  map[string]validator.DDLTable{},
		Operations: map[string]validator.OperationSymbol{},
	}
	funcs := []parser.ServiceFunc{{
		Name: "ListGigs", FileName: "list_gigs.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Gig.List", Inputs: map[string]string{"Query": "query"}, Result: &parser.Result{Type: "Gig", Var: "gigPage", Wrapper: "Page"}},
			{Type: parser.SeqResponse, Target: "gigPage"},
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
