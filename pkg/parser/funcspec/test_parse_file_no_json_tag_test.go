//ff:func feature=funcspec type=parser control=sequence
//ff:what ParseFile JSON 태그 없는 필드 테스트 — json 태그 미지정 시 JSONName이 빈 문자열

package funcspec

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFileNoJSONTag(t *testing.T) {
	dir := t.TempDir()

	src := `package auth

// @func hashPassword
// @description hash

type HashPasswordRequest struct {
	Password string
}

type HashPasswordResponse struct {
	HashedPassword string
}

func HashPassword(req HashPasswordRequest) (HashPasswordResponse, error) {
	return HashPasswordResponse{HashedPassword: "hashed"}, nil
}
`
	path := filepath.Join(dir, "hash_password.go")
	os.WriteFile(path, []byte(src), 0644)

	spec, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}
	if spec == nil {
		t.Fatal("expected non-nil spec")
	}

	// Without json tags, JSONName should be empty.
	if spec.ResponseFields[0].JSONName != "" {
		t.Errorf("JSONName = %q, want empty", spec.ResponseFields[0].JSONName)
	}
}
