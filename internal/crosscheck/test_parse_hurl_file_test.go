//ff:func feature=crosscheck type=rule control=sequence topic=scenario-check
//ff:what TestParseHurlFile: Hurl 파일을 파싱하여 메서드/경로/상태코드 추출 확인
package crosscheck

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseHurlFile(t *testing.T) {
	dir := t.TempDir()
	hurlContent := `# Login
POST {{host}}/auth/login
Content-Type: application/json
{
  "email": "test@test.com",
  "password": "pass"
}

HTTP 200
[Captures]
token: jsonpath "$.access_token"

# Create gig
POST {{host}}/gigs
Authorization: Bearer {{token}}
Content-Type: application/json
{
  "title": "Test"
}

HTTP 201
`

	path := filepath.Join(dir, "scenario-test.hurl")
	os.WriteFile(path, []byte(hurlContent), 0644)

	entries := parseHurlFile(path)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	if entries[0].Method != "POST" || entries[0].Path != "/auth/login" || entries[0].StatusCode != "200" {
		t.Errorf("entry[0] = %+v", entries[0])
	}
	if entries[1].Method != "POST" || entries[1].Path != "/gigs" || entries[1].StatusCode != "201" {
		t.Errorf("entry[1] = %+v", entries[1])
	}
}
