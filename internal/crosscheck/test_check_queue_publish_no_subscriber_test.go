//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=queue-check
//ff:what TestCheckQueuePublishNoSubscriber: publish만 있고 subscribe 없으면 WARNING 확인
package crosscheck

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestCheckQueuePublishNoSubscriber(t *testing.T) {
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
	}

	errs := CheckQueue(funcs, "postgres")

	var found bool
	for _, e := range errs {
		if e.Level == "WARNING" && e.Rule == "Queue publish → subscribe" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected WARNING for publish without subscriber, got: %v", errs)
	}
}
