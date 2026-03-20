//ff:func feature=policy type=parser control=sequence topic=policy-check
//ff:what TestParseFileWithRoles: Rego allow 규칙에서 role 조건 파싱 검증
package policy

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFileWithRoles(t *testing.T) {
	content := `package authz

import rego.v1

default allow := false

allow if {
    input.action == "publish"
    input.resource == "article"
    input.user.role == "editor"
}
`
	dir := t.TempDir()
	path := filepath.Join(dir, "authz.rego")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	p, err := ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(p.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(p.Rules))
	}
	if !p.Rules[0].UsesRole || p.Rules[0].RoleValue != "editor" {
		t.Errorf("expected role=editor, got %+v", p.Rules[0])
	}
}
