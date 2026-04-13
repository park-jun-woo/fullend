//ff:func feature=orchestrator type=util control=iteration dimension=1
//ff:what CWD에서 fullend go.mod를 찾아 pkg/ 경로를 반환
package fullend

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
		if isFullendRoot(dir) {
			return filepath.Join(dir, "pkg")
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}
