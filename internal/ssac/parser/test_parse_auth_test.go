//ff:func feature=ssac-parse type=parser control=sequence
//ff:what TestParseAuth: @auth 인가 어노테이션 파싱 후 액션·리소스·메시지·입력 검증
package parser

import "testing"

func TestParseAuth(t *testing.T) {
	src := `package service

// @auth "delete" "project" {id: project.ID, owner: project.OwnerID} "권한 없음"
func DeleteProject(c *gin.Context) {}
`
	sfs := parseTestFile(t, src)
	seq := sfs[0].Sequences[0]
	assertEqual(t, "Type", seq.Type, SeqAuth)
	assertEqual(t, "Action", seq.Action, "delete")
	assertEqual(t, "Resource", seq.Resource, "project")
	assertEqual(t, "Message", seq.Message, "권한 없음")
	assertEqual(t, "Inputs[id]", seq.Inputs["id"], "project.ID")
	assertEqual(t, "Inputs[owner]", seq.Inputs["owner"], "project.OwnerID")
}
