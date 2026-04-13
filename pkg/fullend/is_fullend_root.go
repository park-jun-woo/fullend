//ff:func feature=orchestrator type=util control=sequence
//ff:what 디렉토리가 fullend go.mod + pkg/ 를 가진 루트인지 판별
package fullend

import (
	"os"
	"path/filepath"
	"strings"
)

// isFullendRoot returns true if dir contains a fullend go.mod and a pkg/ subdirectory.
func isFullendRoot(dir string) bool {
	data, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	if err != nil || !strings.Contains(string(data), "github.com/park-jun-woo/fullend") {
		return false
	}
	fi, err := os.Stat(filepath.Join(dir, "pkg"))
	return err == nil && fi.IsDir()
}
