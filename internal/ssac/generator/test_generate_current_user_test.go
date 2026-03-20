//ff:func feature=ssac-gen type=test control=sequence
//ff:what currentUser 참조 시 MustGet 코드 생성을 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestGenerateCurrentUser(t *testing.T) {
	sf := parser.ServiceFunc{
		Name: "ListMy", FileName: "list_my.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Item.ListByUser", Inputs: map[string]string{"ID": "currentUser.ID"}, Result: &parser.Result{Type: "[]Item", Var: "items"}},
			{Type: parser.SeqResponse, Fields: map[string]string{"items": "items"}},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, `currentUser := c.MustGet("currentUser")`)
	assertContains(t, code, `h.ItemModel.ListByUser(currentUser.ID)`)
}
