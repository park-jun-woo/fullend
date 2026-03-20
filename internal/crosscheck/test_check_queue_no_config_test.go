//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=queue-check
//ff:what TestCheckQueueNoConfig: 큐 설정 없이 publish/subscribe 사용 시 ERROR 확인
package crosscheck

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestCheckQueueNoConfig(t *testing.T) {
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

	errs := CheckQueue(funcs, "")

	var found bool
	for _, e := range errs {
		if e.Level == "ERROR" && e.Rule == "Queue config" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected ERROR for missing queue config, got: %v", errs)
	}
}
