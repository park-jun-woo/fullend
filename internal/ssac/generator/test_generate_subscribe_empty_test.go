//ff:func feature=ssac-gen type=test control=sequence
//ff:what @subscribe에서 @empty가 fmt.Errorf로 생성되는지 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestGenerateSubscribeEmpty(t *testing.T) {
	sf := parser.ServiceFunc{
		Name: "OnTest", FileName: "on_test.go",
		Subscribe: &parser.SubscribeInfo{Topic: "test", MessageType: "TestMsg"},
		Param:     &parser.ParamInfo{TypeName: "TestMsg", VarName: "message"},
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "User.FindByID", Inputs: map[string]string{"ID": "message.UserID"}, Result: &parser.Result{Type: "User", Var: "user"}},
			{Type: parser.SeqEmpty, Target: "user", Message: "사용자 없음"},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, `return fmt.Errorf("사용자 없음")`)
	assertNotContains(t, code, "c.JSON(http.StatusNotFound")
}
