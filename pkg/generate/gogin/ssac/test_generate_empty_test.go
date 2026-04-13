//ff:func feature=ssac-gen type=test control=sequence
//ff:what @empty 가드의 Go 코드 생성을 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateEmpty(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name: "GetCourse", FileName: "get_course.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqGet, Model: "Course.FindByID", Inputs: map[string]string{"ID": "request.ID"}, Result: &ssacparser.Result{Type: "Course", Var: "course"}},
			{Type: ssacparser.SeqEmpty, Target: "course", Message: "코스를 찾을 수 없습니다"},
			{Type: ssacparser.SeqResponse, Fields: map[string]string{"course": "course"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, `if course == nil`)
	assertContains(t, code, `코스를 찾을 수 없습니다`)
}
