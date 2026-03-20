//ff:func feature=ssac-gen type=test control=sequence
//ff:what @subscribe에서 @auth가 fmt.Errorf로 생성되는지 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestGenerateSubscribeAuth(t *testing.T) {
	sf := parser.ServiceFunc{
		Name: "OnTest", FileName: "on_test.go",
		Subscribe: &parser.SubscribeInfo{Topic: "test", MessageType: "TestMsg"},
		Param:     &parser.ParamInfo{TypeName: "TestMsg", VarName: "message"},
		Sequences: []parser.Sequence{
			{Type: parser.SeqAuth, Action: "process", Resource: "order", Inputs: map[string]string{"OrderID": "message.OrderID"}, Message: "Not authorized"},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, `authz.Check(authz.CheckRequest{Action: "process", Resource: "order"`)
	assertContains(t, code, `return fmt.Errorf("Not authorized: %w", err)`)
	assertNotContains(t, code, "c.JSON")
}
