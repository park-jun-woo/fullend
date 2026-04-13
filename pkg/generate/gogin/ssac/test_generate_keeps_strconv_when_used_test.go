//ff:func feature=ssac-gen type=test control=sequence
//ff:what int64 path param이 있을 때 strconv import가 유지되는지 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func TestGenerateKeepsStrconvWhenUsed(t *testing.T) {
	// int64 path param이 있으면 strconv.ParseInt 생성 → strconv 유지
	st := &rule.Ground{
		Models:    map[string]rule.ModelInfo{},
		Tables: map[string]rule.TableInfo{},
		Ops: map[string]rule.OperationInfo{
			"GetCourse": {PathParams: []rule.PathParam{{Name: "ID", GoType: "int64"}}},
		},
	}
	sf := ssacparser.ServiceFunc{
		Name: "GetCourse", FileName: "get_course.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqGet, Model: "Course.FindByID", Inputs: map[string]string{"ID": "request.ID"}, Result: &ssacparser.Result{Type: "Course", Var: "course"}},
			{Type: ssacparser.SeqResponse, Fields: map[string]string{"course": "course"}},
		},
	}
	code := mustGenerate(t, sf, st)
	assertContains(t, code, `"strconv"`)
	assertContains(t, code, `strconv.ParseInt`)
}
