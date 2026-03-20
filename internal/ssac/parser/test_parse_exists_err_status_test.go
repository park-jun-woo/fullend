//ff:func feature=ssac-parse type=parser control=sequence
//ff:what TestParseExistsErrStatus: @exists 커스텀 에러 상태코드(422) 파싱 검증
package parser

import "testing"

func TestParseExistsErrStatus(t *testing.T) {
	src := `package service

// @exists user "Already registered" 422
func Register() {}
`
	sfs := parseTestFile(t, src)
	seq := sfs[0].Sequences[0]
	assertEqual(t, "Type", seq.Type, SeqExists)
	assertEqual(t, "Target", seq.Target, "user")
	assertEqual(t, "Message", seq.Message, "Already registered")
	if seq.ErrStatus != 422 {
		t.Errorf("expected ErrStatus 422, got %d", seq.ErrStatus)
	}
}
