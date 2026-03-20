//ff:func feature=scenario type=parser control=sequence
//ff:what 절대 URL(http/https)이 포함된 .hurl 파일 파싱 검증
package scenario

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFile_AbsoluteURL(t *testing.T) {
	dir := t.TempDir()
	content := `POST http://localhost:8080/auth/register
HTTP 200

GET https://api.example.com/workflows
HTTP 200
`
	path := filepath.Join(dir, "scenario-abs.hurl")
	os.WriteFile(path, []byte(content), 0644)

	entries := ParseFile(path)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	if entries[0].Method != "POST" || entries[0].Path != "/auth/register" {
		t.Errorf("entry[0] = %+v", entries[0])
	}
	if entries[1].Method != "GET" || entries[1].Path != "/workflows" {
		t.Errorf("entry[1] = %+v", entries[1])
	}
}
