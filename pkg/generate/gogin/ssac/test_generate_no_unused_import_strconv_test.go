//ff:func feature=ssac-gen type=test control=sequence
//ff:what strconv import가 불필요할 때 제거되는지 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateNoUnusedImportStrconv(t *testing.T) {
	// string 타입 request param만 있으면 strconv 불필요
	sf := ssacparser.ServiceFunc{
		Name: "GetCourse", FileName: "get_course.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqGet, Model: "Course.FindByID", Inputs: map[string]string{"ID": "request.CourseID"}, Result: &ssacparser.Result{Type: "Course", Var: "course"}},
			{Type: ssacparser.SeqResponse, Fields: map[string]string{"course": "course"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertNotContains(t, code, `"strconv"`)
}
