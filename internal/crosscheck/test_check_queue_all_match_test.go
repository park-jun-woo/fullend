//ff:func feature=crosscheck type=rule control=sequence topic=queue-check
//ff:what TestCheckQueueAllMatch: publish/subscribe 토픽과 필드가 모두 일치하면 에러 없음 확인
package crosscheck

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestCheckQueueAllMatch(t *testing.T) {
	funcs := []ssacparser.ServiceFunc{
		{
			Name: "CreateOrder",
			Sequences: []ssacparser.Sequence{
				{
					Type:   "publish",
					Topic:  "order.completed",
					Inputs: map[string]string{"Email": "order.Email", "OrderID": "order.ID"},
				},
			},
		},
		{
			Name:      "OnOrderCompleted",
			Subscribe: &ssacparser.SubscribeInfo{Topic: "order.completed", MessageType: "OnOrderCompletedMessage"},
			Param:     &ssacparser.ParamInfo{TypeName: "OnOrderCompletedMessage", VarName: "message"},
			Structs: []ssacparser.StructInfo{
				{
					Name: "OnOrderCompletedMessage",
					Fields: []ssacparser.StructField{
						{Name: "Email", Type: "string"},
						{Name: "OrderID", Type: "int64"},
					},
				},
			},
		},
	}

	errs := CheckQueue(funcs, "postgres")

	if len(errs) != 0 {
		t.Errorf("expected no errors, got: %v", errs)
	}
}
