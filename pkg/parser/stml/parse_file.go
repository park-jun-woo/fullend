//ff:func feature=stml-parse type=parser control=sequence
//ff:what 단일 HTML 파일을 파싱하여 PageSpec 반환
package parser

import (
	"os"
	"path/filepath"
)

// ParseFile parses a single HTML file and returns a PageSpec.
func ParseFile(path string) (PageSpec, error) {
	f, err := os.Open(path)
	if err != nil {
		return PageSpec{}, err
	}
	defer f.Close()

	return ParseReader(filepath.Base(path), f)
}
