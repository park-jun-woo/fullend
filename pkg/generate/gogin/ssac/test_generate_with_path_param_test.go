//ff:func feature=ssac-gen type=test control=sequence
//ff:what path parameter가 있을 때 c.Param + strconv 코드 생성을 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func TestGenerateWithPathParam(t *testing.T) {
	st := &rule.Ground{
		Models:    map[string]rule.ModelInfo{},
		Tables: map[string]rule.TableInfo{},
		Ops: map[string]rule.OperationInfo{
			"GetCourse": {
				PathParams: []rule.PathParam{{Name: "CourseID", GoType: "int64"}},
			},
		},
	}
	sf := ssacparser.ServiceFunc{
		Name: "GetCourse", FileName: "get_course.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqGet, Model: "Course.FindByID", Inputs: map[string]string{"CourseID": "request.CourseID"}, Result: &ssacparser.Result{Type: "Course", Var: "course"}},
			{Type: ssacparser.SeqResponse, Fields: map[string]string{"course": "course"}},
		},
	}
	code := mustGenerate(t, sf, st)
	assertContains(t, code, `c.Param("CourseID")`)
	assertContains(t, code, `strconv.ParseInt`)
}
