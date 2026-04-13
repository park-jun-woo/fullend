//ff:func feature=ssac-gen type=test control=sequence
//ff:what @subscribe에서 @auth가 fmt.Errorf로 생성되는지 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateSubscribeAuth(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name: "OnTest", FileName: "on_test.go",
		Subscribe: &ssacparser.SubscribeInfo{Topic: "test", MessageType: "TestMsg"},
		Param:     &ssacparser.ParamInfo{TypeName: "TestMsg", VarName: "message"},
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqAuth, Action: "process", Resource: "order", Inputs: map[string]string{"OrderID": "message.OrderID"}, Message: "Not authorized"},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, `authz.Check(authz.CheckRequest{Action: "process", Resource: "order"`)
	assertContains(t, code, `return fmt.Errorf("Not authorized: %w", err)`)
	assertNotContains(t, code, "c.JSON")
}
