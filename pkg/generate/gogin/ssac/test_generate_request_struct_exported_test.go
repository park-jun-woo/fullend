//ff:func feature=ssac-gen type=test control=sequence
//ff:what request struct 필드가 Exported + json 태그를 갖는지 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateRequestStructExported(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name: "CreateUser", FileName: "create_user.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqPost, Model: "User.Create", Inputs: map[string]string{"email": "request.email"}, Result: &ssacparser.Result{Type: "User", Var: "user"}},
			{Type: ssacparser.SeqResponse, Fields: map[string]string{"user": "user"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, "Email string `json:\"email\"`")
	assertContains(t, code, "email := req.Email")
}
