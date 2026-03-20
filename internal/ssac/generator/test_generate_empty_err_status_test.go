//ff:func feature=ssac-gen type=test control=sequence
//ff:what @empty ErrStatus 커스텀 상태 코드 생성을 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestGenerateEmptyErrStatus(t *testing.T) {
	sf := parser.ServiceFunc{
		Name:     "ActivateWorkflow",
		FileName: "activate_workflow.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqGet, Model: "Org.FindByID", Inputs: map[string]string{"ID": "request.OrgID"}, Result: &parser.Result{Type: "Org", Var: "org"}},
			{Type: parser.SeqEmpty, Target: "org", Message: "Insufficient credits", ErrStatus: 402},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, "http.StatusPaymentRequired")
	assertNotContains(t, code, "http.StatusNotFound")
}
