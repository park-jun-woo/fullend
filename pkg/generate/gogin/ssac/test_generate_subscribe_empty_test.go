//ff:func feature=ssac-gen type=test control=sequence
//ff:what @subscribe에서 @empty가 fmt.Errorf로 생성되는지 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateSubscribeEmpty(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name: "OnTest", FileName: "on_test.go",
		Subscribe: &ssacparser.SubscribeInfo{Topic: "test", MessageType: "TestMsg"},
		Param:     &ssacparser.ParamInfo{TypeName: "TestMsg", VarName: "message"},
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqGet, Model: "User.FindByID", Inputs: map[string]string{"ID": "message.UserID"}, Result: &ssacparser.Result{Type: "User", Var: "user"}},
			{Type: ssacparser.SeqEmpty, Target: "user", Message: "사용자 없음"},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, `return fmt.Errorf("사용자 없음")`)
	assertNotContains(t, code, "c.JSON(http.StatusNotFound")
}
