//ff:func feature=scenario type=parser control=sequence
//ff:what 빈 .hurl 파일에서 항목 0개 반환 검증
package scenario

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFile_Empty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.hurl")
	os.WriteFile(path, []byte("# just a comment\n"), 0644)

	entries := ParseFile(path)
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}
