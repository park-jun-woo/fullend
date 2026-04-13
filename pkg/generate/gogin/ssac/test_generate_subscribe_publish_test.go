//ff:func feature=ssac-gen type=test control=sequence
//ff:what @subscribe에서 @publish의 Go 코드 생성을 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateSubscribePublish(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name: "OnOrderCompleted", FileName: "on_order.go",
		Subscribe: &ssacparser.SubscribeInfo{Topic: "order.completed", MessageType: "Msg"},
		Param:     &ssacparser.ParamInfo{TypeName: "Msg", VarName: "message"},
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqPublish, Topic: "notification.send", Inputs: map[string]string{"Email": "message.Email"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, `queue.Publish(ctx, "notification.send"`)
	assertNotContains(t, code, "c.Request.Context()")
	assertContains(t, code, `"queue"`)
}
