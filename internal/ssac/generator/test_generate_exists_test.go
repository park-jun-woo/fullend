//ff:func feature=ssac-gen type=test control=sequence
//ff:what @exists 가드의 Go 코드 생성을 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestGenerateExists(t *testing.T) {
	sf := parser.ServiceFunc{
		Name: "CreateCourse", FileName: "create_course.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Course.FindByTitle", Inputs: map[string]string{"Title": "request.Title"}, Result: &parser.Result{Type: "Course", Var: "existing"}},
			{Type: parser.SeqExists, Target: "existing", Message: "이미 존재합니다"},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, `if existing != nil`)
	assertContains(t, code, `이미 존재합니다`)
}
