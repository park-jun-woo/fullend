//ff:func feature=policy type=parser control=sequence topic=policy-check
//ff:what TestParseDir: 디렉토리 내 .rego 파일 전체 파싱 검증
package policy

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseDir(t *testing.T) {
	dir := t.TempDir()
	content := `package authz

import rego.v1

default allow := false

allow if {
    input.action == "create"
    input.resource == "item"
}
`
	if err := os.WriteFile(filepath.Join(dir, "authz.rego"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	policies, err := ParseDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(policies) != 1 {
		t.Fatalf("expected 1 policy, got %d", len(policies))
	}
	if len(policies[0].Rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(policies[0].Rules))
	}
}
