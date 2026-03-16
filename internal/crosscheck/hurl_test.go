package crosscheck

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestNormalizeHurlPath(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"/gigs/{{gig_id}}", []string{"gigs", ":param"}},
		{"/gigs", []string{"gigs"}},
		{"/gigs/{{gig_id}}/proposals", []string{"gigs", ":param", "proposals"}},
		{"/auth/login", []string{"auth", "login"}},
		{"/gigs?status=open", []string{"gigs"}},
	}

	for _, tt := range tests {
		got := normalizeHurlPath(tt.input)
		if !slices.Equal(got, tt.want) {
			t.Errorf("normalizeHurlPath(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestNormalizeOpenAPIPath(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"/gigs/{id}", []string{"gigs", ":param"}},
		{"/gigs/{gigId}/proposals/{proposalId}", []string{"gigs", ":param", "proposals", ":param"}},
		{"/auth/login", []string{"auth", "login"}},
	}

	for _, tt := range tests {
		got := normalizeOpenAPIPath(tt.input)
		if !slices.Equal(got, tt.want) {
			t.Errorf("normalizeOpenAPIPath(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestSegmentsMatch(t *testing.T) {
	tests := []struct {
		a, b []string
		want bool
	}{
		{[]string{"gigs", ":param"}, []string{"gigs", ":param"}, true},
		{[]string{"gigs"}, []string{"gigs", ":param"}, false},
		{[]string{"gigs", ":param"}, []string{"users", ":param"}, false},
		{[]string{"auth", "login"}, []string{"auth", "login"}, true},
	}

	for _, tt := range tests {
		got := segmentsMatch(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("segmentsMatch(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

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
