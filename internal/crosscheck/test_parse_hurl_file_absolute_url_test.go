//ff:func feature=crosscheck type=rule control=sequence topic=scenario-check
//ff:what TestParseHurlFile_AbsoluteURL: 절대 URL이 포함된 Hurl 파일 파싱 시 경로만 추출되는지 확인
package crosscheck

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseHurlFile_AbsoluteURL(t *testing.T) {
	dir := t.TempDir()
	hurlContent := `# Register
POST http://localhost:8080/auth/register
Content-Type: application/json
{
  "email": "test@test.com",
  "password": "pass"
}

HTTP 200

# Login
POST http://localhost:8080/auth/login
Content-Type: application/json
{
  "email": "test@test.com",
  "password": "pass"
}

HTTP 200

# Create with HTTPS
GET https://api.example.com/workflows
Authorization: Bearer {{token}}

HTTP 200
`

	path := filepath.Join(dir, "scenario-abs.hurl")
	os.WriteFile(path, []byte(hurlContent), 0644)

	entries := parseHurlFile(path)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}

	if entries[0].Method != "POST" || entries[0].Path != "/auth/register" || entries[0].StatusCode != "200" {
		t.Errorf("entry[0] = %+v", entries[0])
	}
	if entries[1].Method != "POST" || entries[1].Path != "/auth/login" || entries[1].StatusCode != "200" {
		t.Errorf("entry[1] = %+v", entries[1])
	}
	if entries[2].Method != "GET" || entries[2].Path != "/workflows" || entries[2].StatusCode != "200" {
		t.Errorf("entry[2] = %+v", entries[2])
	}
}
