//ff:func feature=ssac-parse type=parser control=sequence
//ff:what TestParseAuthErrStatus: @auth 커스텀 에러 상태코드(401) 파싱 검증
package parser

import "testing"

func TestParseAuthErrStatus(t *testing.T) {
	src := `package service

// @auth "ActivateWorkflow" "workflow" {UserID: currentUser.ID, ResourceID: workflow.ID} "Forbidden" 401
func ActivateWorkflow() {}
`
	sfs := parseTestFile(t, src)
	seq := sfs[0].Sequences[0]
	assertEqual(t, "Type", seq.Type, SeqAuth)
	assertEqual(t, "Action", seq.Action, "ActivateWorkflow")
	assertEqual(t, "Resource", seq.Resource, "workflow")
	assertEqual(t, "Message", seq.Message, "Forbidden")
	if seq.ErrStatus != 401 {
		t.Errorf("expected ErrStatus 401, got %d", seq.ErrStatus)
	}
}
