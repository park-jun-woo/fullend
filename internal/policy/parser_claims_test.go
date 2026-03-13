package policy

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFile_ClaimsRefs(t *testing.T) {
	dir := t.TempDir()
	regoContent := `package authz

default allow = false

allow if {
    input.action == "create"
    input.resource == "gig"
    input.claims.role == "client"
    data.owners.gig[input.resource_id] == input.claims.user_id
}
`
	path := filepath.Join(dir, "authz.rego")
	if err := os.WriteFile(path, []byte(regoContent), 0644); err != nil {
		t.Fatal(err)
	}

	p, err := ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}

	// Should have role and user_id as claims refs.
	refs := make(map[string]bool)
	for _, r := range p.ClaimsRefs {
		refs[r] = true
	}

	if !refs["role"] {
		t.Error("expected ClaimsRefs to contain 'role'")
	}
	if !refs["user_id"] {
		t.Error("expected ClaimsRefs to contain 'user_id'")
	}
	if len(p.ClaimsRefs) != 2 {
		t.Errorf("expected 2 unique claims refs, got %d: %v", len(p.ClaimsRefs), p.ClaimsRefs)
	}
}
