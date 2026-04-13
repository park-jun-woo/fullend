//ff:func feature=ssac-gen type=test control=sequence
//ff:what @subscribe 함수에서 @auth claims가 포함되는지 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateSubscribeAuthClaims(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name: "OnTest", FileName: "on_test.go",
		Subscribe: &ssacparser.SubscribeInfo{Topic: "test", MessageType: "TestMsg"},
		Param:     &ssacparser.ParamInfo{TypeName: "TestMsg", VarName: "message"},
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqAuth, Action: "process", Resource: "order", Inputs: map[string]string{"UserID": "currentUser.ID"}, Message: "Not authorized"},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, `UserID: currentUser.ID`)
}
