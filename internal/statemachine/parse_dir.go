//ff:func feature=statemachine type=parser control=iteration dimension=1 topic=states
//ff:what 디렉토리 내 모든 .md 파일을 파싱하여 StateDiagram 슬라이스를 반환한다
package statemachine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ParseDir parses all *.md files in the given directory.
func ParseDir(dir string) ([]*StateDiagram, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read states dir: %w", err)
	}

	var diagrams []*StateDiagram
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		d, err := ParseFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		diagrams = append(diagrams, d)
	}
	return diagrams, nil
}
