//ff:func feature=ssac-gen type=test control=sequence
//ff:what path parameter가 있을 때 c.Param + strconv 코드 생성을 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func TestGenerateWithPathParam(t *testing.T) {
	st := &validator.SymbolTable{
		Models:    map[string]validator.ModelSymbol{},
		DDLTables: map[string]validator.DDLTable{},
		Operations: map[string]validator.OperationSymbol{
			"GetCourse": {
				PathParams:    []validator.PathParam{{Name: "CourseID", GoType: "int64"}},
				RequestFields: map[string]bool{"CourseID": true},
			},
		},
	}
	sf := parser.ServiceFunc{
		Name: "GetCourse", FileName: "get_course.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Course.FindByID", Inputs: map[string]string{"CourseID": "request.CourseID"}, Result: &parser.Result{Type: "Course", Var: "course"}},
			{Type: parser.SeqResponse, Fields: map[string]string{"course": "course"}},
		},
	}
	code := mustGenerate(t, sf, st)
	assertContains(t, code, `c.Param("CourseID")`)
	assertContains(t, code, `strconv.ParseInt`)
}
