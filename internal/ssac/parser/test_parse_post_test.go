//ff:func feature=ssac-parse type=parser control=sequence
//ff:what TestParsePost: @post 어노테이션 파싱 후 타입·결과·입력 검증
package parser

import "testing"

func TestParsePost(t *testing.T) {
	src := `package service

// @post Session session = Session.Create({ProjectID: request.ProjectID, Command: request.Command})
func CreateSession(c *gin.Context) {}
`
	sfs := parseTestFile(t, src)
	seq := sfs[0].Sequences[0]
	assertEqual(t, "Type", seq.Type, SeqPost)
	assertEqual(t, "Result.Type", seq.Result.Type, "Session")
	assertEqual(t, "Result.Var", seq.Result.Var, "session")
	if len(seq.Inputs) != 2 {
		t.Fatalf("expected 2 inputs, got %d", len(seq.Inputs))
	}
}
