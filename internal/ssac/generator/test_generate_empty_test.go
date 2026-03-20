//ff:func feature=ssac-gen type=test control=sequence
//ff:what @empty 가드의 Go 코드 생성을 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestGenerateEmpty(t *testing.T) {
	sf := parser.ServiceFunc{
		Name: "GetCourse", FileName: "get_course.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Course.FindByID", Inputs: map[string]string{"ID": "request.ID"}, Result: &parser.Result{Type: "Course", Var: "course"}},
			{Type: parser.SeqEmpty, Target: "course", Message: "코스를 찾을 수 없습니다"},
			{Type: parser.SeqResponse, Fields: map[string]string{"course": "course"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, `if course == nil`)
	assertContains(t, code, `코스를 찾을 수 없습니다`)
}
