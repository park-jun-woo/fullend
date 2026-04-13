//ff:func feature=ssac-gen type=test control=sequence
//ff:what @exists 가드의 Go 코드 생성을 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateExists(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name: "CreateCourse", FileName: "create_course.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqGet, Model: "Course.FindByTitle", Inputs: map[string]string{"Title": "request.Title"}, Result: &ssacparser.Result{Type: "Course", Var: "existing"}},
			{Type: ssacparser.SeqExists, Target: "existing", Message: "이미 존재합니다"},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, `if existing != nil`)
	assertContains(t, code, `이미 존재합니다`)
}
