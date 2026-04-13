//ff:func feature=ssac-gen type=test control=sequence
//ff:what @get 시퀀스의 Go 코드 생성을 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateGet(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name: "GetCourse", FileName: "get_course.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqGet, Model: "Course.FindByID", Inputs: map[string]string{"CourseID": "request.CourseID"}, Result: &ssacparser.Result{Type: "Course", Var: "course"}},
			{Type: ssacparser.SeqResponse, Fields: map[string]string{"course": "course"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, `course, err := h.CourseModel.FindByID(courseID)`)
	assertContains(t, code, `courseID := c.Query("CourseID")`)
	assertContains(t, code, `"course": course`)
}
