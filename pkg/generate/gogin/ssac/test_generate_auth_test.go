//ff:func feature=ssac-gen type=test control=sequence
//ff:what @auth 가드의 Go 코드 생성을 검증
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestGenerateAuth(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name: "DeleteProject", FileName: "delete_project.go",
		Sequences: []ssacparser.Sequence{
			{Type: ssacparser.SeqAuth, Action: "delete", Resource: "project", Inputs: map[string]string{"id": "project.ID", "owner": "project.OwnerID"}, Message: "권한 없음"},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, `authz.Check(authz.CheckRequest{Action: "delete", Resource: "project"`)
	assertContains(t, code, `ID: project.ID`)
	assertContains(t, code, `Owner: project.OwnerID`)
	assertContains(t, code, `http.StatusForbidden`)
	assertNotContains(t, code, `authz.Input{`)
}
