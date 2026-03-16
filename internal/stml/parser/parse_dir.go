//ff:func feature=stml-parse type=parser control=iteration dimension=1
//ff:what 디렉토리 내 모든 .html 파일을 파싱하여 PageSpec 목록 반환
package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ParseDir parses all .html files in the given directory and returns a PageSpec for each.
func ParseDir(dir string) ([]PageSpec, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read dir %s: %w", dir, err)
	}

	var pages []PageSpec
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".html") {
			continue
		}
		page, err := ParseFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", e.Name(), err)
		}
		pages = append(pages, page)
	}
	return pages, nil
}
