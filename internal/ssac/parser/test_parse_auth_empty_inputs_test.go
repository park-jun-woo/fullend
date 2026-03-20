//ff:func feature=ssac-parse type=parser control=sequence
//ff:what TestParseAuthEmptyInputs: @auth 빈 입력({}) 파싱 검증
package parser

import "testing"

func TestParseAuthEmptyInputs(t *testing.T) {
	src := `package service

// @auth "view" "dashboard" {} "권한 없음"
func ViewDashboard(c *gin.Context) {}
`
	sfs := parseTestFile(t, src)
	seq := sfs[0].Sequences[0]
	assertEqual(t, "Action", seq.Action, "view")
	if len(seq.Inputs) != 0 {
		t.Errorf("expected empty inputs, got %d", len(seq.Inputs))
	}
}
