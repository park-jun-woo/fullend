//ff:func feature=ssac-gen type=test control=sequence
//ff:what @empty ErrStatus 커스텀 상태 코드 생성을 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateEmptyErrStatus(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name:     "ActivateWorkflow",
		FileName: "activate_workflow.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqGet, Model: "Org.FindByID", Inputs: map[string]string{"ID": "request.OrgID"}, Result: &ssacparser.Result{Type: "Org", Var: "org"}},
			{Type: ssacparser.SeqEmpty, Target: "org", Message: "Insufficient credits", ErrStatus: 402},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, "http.StatusPaymentRequired")
	assertNotContains(t, code, "http.StatusNotFound")
}
