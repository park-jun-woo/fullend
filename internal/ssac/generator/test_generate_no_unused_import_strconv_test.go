//ff:func feature=ssac-gen type=test control=sequence
//ff:what strconv import가 불필요할 때 제거되는지 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestGenerateNoUnusedImportStrconv(t *testing.T) {
	// string 타입 request param만 있으면 strconv 불필요
	sf := parser.ServiceFunc{
		Name: "GetCourse", FileName: "get_course.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Course.FindByID", Inputs: map[string]string{"ID": "request.CourseID"}, Result: &parser.Result{Type: "Course", Var: "course"}},
			{Type: parser.SeqResponse, Fields: map[string]string{"course": "course"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertNotContains(t, code, `"strconv"`)
}
