//ff:func feature=ssac-gen type=test control=sequence
//ff:what @response 시퀀스의 Go 코드 생성을 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateResponse(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name: "GetCourse", FileName: "get_course.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqGet, Model: "Course.FindByID", Inputs: map[string]string{"ID": "request.ID"}, Result: &ssacparser.Result{Type: "Course", Var: "course"}},
			{Type: ssacparser.SeqGet, Model: "User.FindByID", Inputs: map[string]string{"InstructorID": "course.InstructorID"}, Result: &ssacparser.Result{Type: "User", Var: "instructor"}},
			{Type: ssacparser.SeqResponse, Fields: map[string]string{"course": "course", "instructor_name": "instructor.Name"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, `"course":`)
	assertContains(t, code, `"instructor_name":`)
	assertContains(t, code, `instructor.Name`)
}
