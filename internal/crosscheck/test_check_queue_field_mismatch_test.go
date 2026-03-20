//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=queue-check
//ff:what TestCheckQueueFieldMismatch: publish 입력과 subscribe 메시지 필드 불일치 시 WARNING 확인
package crosscheck

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestCheckQueueFieldMismatch(t *testing.T) {
	funcs := []ssacparser.ServiceFunc{
		{
			Name: "CreateOrder",
			Sequences: []ssacparser.Sequence{
				{
					Type:   "publish",
					Topic:  "order.completed",
					Inputs: map[string]string{"Email": "order.Email"},
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

	var found bool
	for _, e := range errs {
		if e.Level == "WARNING" && e.Rule == "Queue field mismatch" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected WARNING for field mismatch, got: %v", errs)
	}
}
