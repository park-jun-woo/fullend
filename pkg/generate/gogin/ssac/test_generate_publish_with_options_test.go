//ff:func feature=ssac-gen type=test control=sequence
//ff:what @publish options(delay) 지정 시 WithDelay 코드 생성을 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGeneratePublishWithOptions(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name: "AbandonCart", FileName: "abandon_cart.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqPublish, Topic: "cart.abandoned", Inputs: map[string]string{"CartID": "cart.ID"}, Options: map[string]string{"delay": "1800"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, `queue.Publish(c.Request.Context(), "cart.abandoned"`)
	assertContains(t, code, `queue.WithDelay(1800)`)
}
