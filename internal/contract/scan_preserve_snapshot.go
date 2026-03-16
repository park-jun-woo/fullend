//ff:func feature=contract type=walker control=sequence
//ff:what 디렉토리를 순회하여 보존 대상 콘텐츠를 캡처한다
package contract

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ScanPreserveSnapshot walks a directory and captures all preserved content.
func ScanPreserveSnapshot(dir string) *PreserveSnapshot {
	snap := &PreserveSnapshot{
		FilePreserves: make(map[string]string),
		FuncPreserves: make(map[string]map[string]*PreservedFunc),
	}

	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		src, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		content := string(src)

		// Check file-level preserve.
		if hasFilePreserve(content) {
			snap.FilePreserves[path] = content
			return nil
		}

		// Check function-level preserves.
		funcs := scanPreservedFromSource(content)
		if len(funcs) > 0 {
			snap.FuncPreserves[path] = funcs
		}
		return nil
	})

	return snap
}
