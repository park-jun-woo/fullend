//ff:func feature=ssac-parse type=parser control=sequence
//ff:what 콜론 없는 입력 형식 에러 검증 — {query} 형태 거부

package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseInputsNoColon(t *testing.T) {
	src := `package service

// @get []Gig gigs = Gig.List({query})
func ListGigs(c *gin.Context) {}
`
	dir := t.TempDir()
	path := filepath.Join(dir, "test.go")
	if err := os.WriteFile(path, []byte(src), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := ParseFile(path)
	if err == nil {
		t.Fatal("expected error for input without colon")
	}
	if !strings.Contains(err.Error(), "유효하지 않은 입력 형식") {
		t.Errorf("unexpected error: %v", err)
	}
}
