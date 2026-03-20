//ff:func feature=orchestrator type=util control=iteration dimension=1
//ff:what findSpecsDir: 프로젝트 루트에서 specs/gigbridge 디렉토리 경로 탐색 헬퍼
package orchestrator

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// findSpecsDir locates specs/gigbridge relative to the project root.
func findSpecsDir(t *testing.T) string {
	t.Helper()
	// Walk up from this test file to find project root (where go.mod lives).
	_, thisFile, _, _ := runtime.Caller(0)
	dir := filepath.Dir(thisFile)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("cannot find project root (go.mod)")
		}
		dir = parent
	}
	specsDir := filepath.Join(dir, "specs", "gigbridge")
	if _, err := os.Stat(specsDir); err != nil {
		t.Skipf("specs/gigbridge not found: %v", err)
	}
	return specsDir
}
