//ff:func feature=ssac-gen type=test control=sequence
//ff:what 모델 인터페이스 Go 코드 생성을 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func TestGenerateModelInterface(t *testing.T) {
	st := &validator.SymbolTable{
		Models: map[string]validator.ModelSymbol{
			"Course": {Methods: map[string]validator.MethodInfo{
				"FindByID": {Cardinality: "one"},
			}},
		},
		DDLTables: map[string]validator.DDLTable{
			"courses": {Columns: map[string]string{"id": "int64", "title": "string"}},
		},
		Operations: map[string]validator.OperationSymbol{},
	}
	funcs := []ssacparser.ServiceFunc{{
		Name: "GetCourse", FileName: "get_course.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqGet, Model: "Course.FindByID", Inputs: map[string]string{"CourseID": "request.CourseID"}, Result: &ssacparser.Result{Type: "Course", Var: "course"}},
			{Type: ssacparser.SeqResponse, Fields: map[string]string{"course": "course"}},
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
	assertContains(t, data, "type CourseModel interface")
	assertContains(t, data, "WithTx(tx *sql.Tx) CourseModel")
	assertContains(t, data, "FindByID(")
}
