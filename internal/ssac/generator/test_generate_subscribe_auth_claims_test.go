//ff:func feature=ssac-gen type=test control=sequence
//ff:what @subscribe 함수에서 @auth claims가 포함되는지 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestGenerateSubscribeAuthClaims(t *testing.T) {
	sf := parser.ServiceFunc{
		Name: "OnTest", FileName: "on_test.go",
		Subscribe: &parser.SubscribeInfo{Topic: "test", MessageType: "TestMsg"},
		Param:     &parser.ParamInfo{TypeName: "TestMsg", VarName: "message"},
		Sequences: []parser.Sequence{
			{Type: parser.SeqAuth, Action: "process", Resource: "order", Inputs: map[string]string{"UserID": "currentUser.ID"}, Message: "Not authorized"},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, `UserID: currentUser.ID`)
}
