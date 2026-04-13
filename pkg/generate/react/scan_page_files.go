//ff:func feature=gen-react type=util control=iteration dimension=1
//ff:what pages 디렉토리에서 .tsx 파일 목록을 스캔한다

package react

import (
	"os"
	"sort"
	"strings"
)

// scanPageFiles scans a directory for .tsx files and returns sorted file names without extension.
func scanPageFiles(pagesDir string) []string {
	var pageFiles []string
	entries, err := os.ReadDir(pagesDir)
	if err != nil {
		return nil
	}
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".tsx") {
			pageFiles = append(pageFiles, strings.TrimSuffix(e.Name(), ".tsx"))
		}
	}
	sort.Strings(pageFiles)
	return pageFiles
}
