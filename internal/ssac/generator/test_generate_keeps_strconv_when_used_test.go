//ff:func feature=ssac-gen type=test control=sequence
//ff:what int64 path param이 있을 때 strconv import가 유지되는지 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func TestGenerateKeepsStrconvWhenUsed(t *testing.T) {
	// int64 path param이 있으면 strconv.ParseInt 생성 → strconv 유지
	st := &validator.SymbolTable{
		Models:    map[string]validator.ModelSymbol{},
		DDLTables: map[string]validator.DDLTable{},
		Operations: map[string]validator.OperationSymbol{
			"GetCourse": {PathParams: []validator.PathParam{{Name: "ID", GoType: "int64"}}},
		},
	}
	sf := parser.ServiceFunc{
		Name: "GetCourse", FileName: "get_course.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Course.FindByID", Inputs: map[string]string{"ID": "request.ID"}, Result: &parser.Result{Type: "Course", Var: "course"}},
			{Type: parser.SeqResponse, Fields: map[string]string{"course": "course"}},
		},
	}
	code := mustGenerate(t, sf, st)
	assertContains(t, code, `"strconv"`)
	assertContains(t, code, `strconv.ParseInt`)
}
