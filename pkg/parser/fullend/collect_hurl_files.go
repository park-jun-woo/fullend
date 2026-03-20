//ff:func feature=orchestrator type=util control=iteration dimension=1
//ff:what 디렉토리에서 .hurl 파일 경로 목록을 수집
package fullend

import (
	"os"
	"path/filepath"
	"strings"
)

// collectHurlFiles returns all .hurl file paths in the given directory.
func collectHurlFiles(dir string) []string {
	var files []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".hurl") {
			files = append(files, filepath.Join(dir, e.Name()))
		}
	}
	return files
}
