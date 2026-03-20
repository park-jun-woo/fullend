//ff:func feature=ssac-parse type=parser control=sequence
//ff:what parseTestFile 헬퍼 — 소스 문자열을 임시 파일에 쓰고 ParseFile 호출

package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func parseTestFile(t *testing.T, src string) []ServiceFunc {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.go")
	if err := os.WriteFile(path, []byte(src), 0644); err != nil {
		t.Fatal(err)
	}
	sfs, err := ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(sfs) == 0 {
		t.Fatal("expected at least 1 ServiceFunc")
	}
	return sfs
}
