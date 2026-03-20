//ff:func feature=ssac-gen type=test control=sequence
//ff:what @subscribe에서 @publish의 Go 코드 생성을 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestGenerateSubscribePublish(t *testing.T) {
	sf := parser.ServiceFunc{
		Name: "OnOrderCompleted", FileName: "on_order.go",
		Subscribe: &parser.SubscribeInfo{Topic: "order.completed", MessageType: "Msg"},
		Param:     &parser.ParamInfo{TypeName: "Msg", VarName: "message"},
		Sequences: []parser.Sequence{
			{Type: parser.SeqPublish, Topic: "notification.send", Inputs: map[string]string{"Email": "message.Email"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, `queue.Publish(ctx, "notification.send"`)
	assertNotContains(t, code, "c.Request.Context()")
	assertContains(t, code, `"queue"`)
}
