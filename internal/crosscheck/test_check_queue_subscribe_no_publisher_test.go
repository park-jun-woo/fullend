//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=queue-check
//ff:what TestCheckQueueSubscribeNoPublisher: subscribe만 있고 publish 없으면 WARNING 확인
package crosscheck

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestCheckQueueSubscribeNoPublisher(t *testing.T) {
	funcs := []ssacparser.ServiceFunc{
		{
			Name:      "OnOrderCompleted",
			Subscribe: &ssacparser.SubscribeInfo{Topic: "order.completed", MessageType: "OnOrderCompletedMessage"},
			Param:     &ssacparser.ParamInfo{TypeName: "OnOrderCompletedMessage", VarName: "message"},
			Structs: []ssacparser.StructInfo{
				{
					Name:   "OnOrderCompletedMessage",
					Fields: []ssacparser.StructField{{Name: "Email", Type: "string"}},
				},
			},
		},
	}

	errs := CheckQueue(funcs, "postgres")

	var found bool
	for _, e := range errs {
		if e.Level == "WARNING" && e.Rule == "Queue subscribe → publish" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected WARNING for subscribe without publisher, got: %v", errs)
	}
}
