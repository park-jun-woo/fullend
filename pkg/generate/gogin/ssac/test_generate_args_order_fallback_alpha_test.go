//ff:func feature=ssac-gen type=test control=sequence
//ff:what 심볼 테이블 없을 때 알파벳순 fallback 인자 순서 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateArgsOrderFallbackAlpha(t *testing.T) {
	// 심볼 테이블 없으면 알파벳순 fallback
	sf := ssacparser.ServiceFunc{
		Name: "PublishGig", FileName: "publish_gig.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqPut, Model: "Gig.UpdateStatus", Inputs: map[string]string{"ID": "request.GigID", "Status": `"published"`}},
		},
	}
	code := mustGenerate(t, sf, nil)
	// 알파벳순: ID < Status → (gigID, "published")
	assertContains(t, code, `h.GigModel.WithTx(tx).UpdateStatus(gigID, "published")`)
}
