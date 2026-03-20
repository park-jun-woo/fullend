//ff:func feature=ssac-gen type=test control=sequence
//ff:what @get 시퀀스의 Go 코드 생성을 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestGenerateGet(t *testing.T) {
	sf := parser.ServiceFunc{
		Name: "GetCourse", FileName: "get_course.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Course.FindByID", Inputs: map[string]string{"CourseID": "request.CourseID"}, Result: &parser.Result{Type: "Course", Var: "course"}},
			{Type: parser.SeqResponse, Fields: map[string]string{"course": "course"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, `course, err := h.CourseModel.FindByID(courseID)`)
	assertContains(t, code, `courseID := c.Query("CourseID")`)
	assertContains(t, code, `"course": course`)
}
