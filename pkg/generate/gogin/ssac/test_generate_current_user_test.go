//ff:func feature=ssac-gen type=test control=sequence
//ff:what currentUser 참조 시 MustGet 코드 생성을 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateCurrentUser(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name: "ListMy", FileName: "list_my.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqGet, Model: "Item.ListByUser", Inputs: map[string]string{"ID": "currentUser.ID"}, Result: &ssacparser.Result{Type: "[]Item", Var: "items"}},
			{Type: ssacparser.SeqResponse, Fields: map[string]string{"items": "items"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, `currentUser := c.MustGet("currentUser")`)
	assertContains(t, code, `h.ItemModel.ListByUser(currentUser.ID)`)
}
