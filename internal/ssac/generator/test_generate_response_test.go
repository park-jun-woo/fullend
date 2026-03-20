//ff:func feature=ssac-gen type=test control=sequence
//ff:what @response 시퀀스의 Go 코드 생성을 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestGenerateResponse(t *testing.T) {
	sf := parser.ServiceFunc{
		Name: "GetCourse", FileName: "get_course.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Course.FindByID", Inputs: map[string]string{"ID": "request.ID"}, Result: &parser.Result{Type: "Course", Var: "course"}},
			{Type: parser.SeqGet, Model: "User.FindByID", Inputs: map[string]string{"InstructorID": "course.InstructorID"}, Result: &parser.Result{Type: "User", Var: "instructor"}},
			{Type: parser.SeqResponse, Fields: map[string]string{"course": "course", "instructor_name": "instructor.Name"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, `"course":`)
	assertContains(t, code, `"instructor_name":`)
	assertContains(t, code, `instructor.Name`)
}
