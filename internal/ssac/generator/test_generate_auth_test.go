//ff:func feature=ssac-gen type=test control=sequence
//ff:what @auth 가드의 Go 코드 생성을 검증
package generator

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestGenerateAuth(t *testing.T) {
	sf := parser.ServiceFunc{
		Name: "DeleteProject", FileName: "delete_project.go",
		Sequences: []parser.Sequence{
			{Type: parser.SeqAuth, Action: "delete", Resource: "project", Inputs: map[string]string{"id": "project.ID", "owner": "project.OwnerID"}, Message: "권한 없음"},
		},
	}
	code := mustGenerate(t, sf, nil)
	assertContains(t, code, `authz.Check(authz.CheckRequest{Action: "delete", Resource: "project"`)
	assertContains(t, code, `ID: project.ID`)
	assertContains(t, code, `Owner: project.OwnerID`)
	assertContains(t, code, `http.StatusForbidden`)
	assertNotContains(t, code, `authz.Input{`)
}
