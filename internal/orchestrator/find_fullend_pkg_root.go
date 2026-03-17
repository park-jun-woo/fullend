//ff:func feature=orchestrator type=util control=iteration dimension=2
//ff:what fullend pkg/ 디렉토리 탐색 — CWD에서 상위로 go.mod 검색
package orchestrator

import (
	"os"
	"path/filepath"
)

// findFullendPkgRoot locates the fullend pkg/ directory.
// Walks up from CWD looking for go.mod with module github.com/park-jun-woo/fullend.
func findFullendPkgRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	for {
		goModPath := filepath.Join(dir, "go.mod")
		if data, err := os.ReadFile(goModPath); err == nil && isFullendGoMod(data) {
			pkgDir := filepath.Join(dir, "pkg")
			if fi, err := os.Stat(pkgDir); err == nil && fi.IsDir() {
				return pkgDir
			}
			return ""
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}
