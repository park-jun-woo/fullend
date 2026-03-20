package scenario

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFile(t *testing.T) {
	dir := t.TempDir()
	content := `# Login
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
	os.WriteFile(path, []byte(content), 0644)

	entries := ParseFile(path)
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

func TestParseFile_Empty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.hurl")
	os.WriteFile(path, []byte("# just a comment\n"), 0644)

	entries := ParseFile(path)
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestParseFile_NotFound(t *testing.T) {
	entries := ParseFile("/nonexistent/file.hurl")
	if entries != nil {
		t.Fatalf("expected nil, got %v", entries)
	}
}
