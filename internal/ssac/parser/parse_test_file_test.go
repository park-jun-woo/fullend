//ff:func feature=ssac-parse type=util control=sequence
//ff:what parseTestFile: 테스트용 SSaC 소스를 임시 파일에 쓰고 ParseFile 호출하는 헬퍼
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
