//ff:func feature=ssac-gen type=test control=sequence
//ff:what @publish options(delay) 지정 시 WithDelay 코드 생성을 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestGeneratePublishWithOptions(t *testing.T) {
	sf := parser.ServiceFunc{
		Name: "AbandonCart", FileName: "abandon_cart.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqPublish, Topic: "cart.abandoned", Inputs: map[string]string{"CartID": "cart.ID"}, Options: map[string]string{"delay": "1800"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, `queue.Publish(c.Request.Context(), "cart.abandoned"`)
	assertContains(t, code, `queue.WithDelay(1800)`)
}
