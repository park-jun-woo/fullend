//ff:func feature=contract type=walker control=sequence
//ff:what 디렉토리 내 보존된 함수 및 파일 수를 센다
package contract

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// CountPreserveFuncs counts all preserved functions and files in a directory.
func CountPreserveFuncs(dir string) int {
	count := 0
	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		src, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		content := string(src)
		if hasFilePreserve(content) {
			count++
			return nil
		}
		funcs := scanPreservedFromSource(content)
		count += len(funcs)
		return nil
	})
	return count
}
