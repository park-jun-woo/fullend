//ff:func feature=funcspec type=parser control=sequence
//ff:what 디렉토리 내 모든 .go 파일을 재귀적으로 파싱하여 FuncSpec 목록을 반환한다
package funcspec

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

// ParseDir parses all .go files under dir (recursively by package subdirectory).
// Returns a flat list of FuncSpecs.
func ParseDir(dir string) ([]FuncSpec, error) {
	var specs []FuncSpec
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(d.Name(), ".go") {
			return err
		}
		fs, err := ParseFile(path)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		if fs != nil {
			// Derive package from parent dir name.
			rel, _ := filepath.Rel(dir, path)
			parts := strings.Split(filepath.Dir(rel), string(filepath.Separator))
			if parts[0] != "." {
				fs.Package = parts[0]
			}
			specs = append(specs, *fs)
		}
		return nil
	})
	return specs, err
}
