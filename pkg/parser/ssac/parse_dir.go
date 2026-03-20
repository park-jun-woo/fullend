//ff:func feature=ssac-parse type=parser control=sequence
//ff:what 디렉토리 내 모든 .ssac 파일을 재귀 탐색하여 []ServiceFunc 반환
package parser

import (
	"io/fs"
	"path/filepath"
	"strings"
)

// ParseDir은 디렉토리 내 모든 .ssac 파일을 재귀 탐색하여 []ServiceFunc를 반환한다.
func ParseDir(dir string) ([]ServiceFunc, error) {
	var funcs []ServiceFunc
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(d.Name(), ".ssac") {
			return err
		}
		parsed, walkErr := parseDirEntry(dir, path, d.Name())
		if walkErr != nil {
			return walkErr
		}
		funcs = append(funcs, parsed...)
		return nil
	})
	return funcs, err
}
