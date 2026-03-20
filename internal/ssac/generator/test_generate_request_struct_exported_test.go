//ff:func feature=ssac-gen type=test control=sequence
//ff:what request struct 필드가 Exported + json 태그를 갖는지 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestGenerateRequestStructExported(t *testing.T) {
	sf := parser.ServiceFunc{
		Name: "CreateUser", FileName: "create_user.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqPost, Model: "User.Create", Inputs: map[string]string{"email": "request.email"}, Result: &parser.Result{Type: "User", Var: "user"}},
			{Type: parser.SeqResponse, Fields: map[string]string{"user": "user"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, "Email string `json:\"email\"`")
	assertContains(t, code, "email := req.Email")
}
