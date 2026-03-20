//ff:func feature=statemachine type=parser control=sequence topic=states
//ff:what Mermaid stateDiagram 마크다운 파일 하나를 파싱한다
package statemachine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ParseFile parses a single Mermaid stateDiagram markdown file.
// The diagram ID is derived from the filename (without extension).
func ParseFile(path string) (*StateDiagram, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read state file %s: %w", path, err)
	}

	id := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	return Parse(id, string(data))
}
