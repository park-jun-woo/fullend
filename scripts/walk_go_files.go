//ff:func feature=stat type=util control=sequence
//ff:what 디렉토리를 순회하며 .go 파일마다 콜백 호출
package main

import (
	"os"
	"path/filepath"
	"strings"
)

func walkGoFiles(dir string, fn func(path string)) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		if strings.Contains(path, "vendor/") || strings.Contains(path, "_test.go") {
			return nil
		}
		fn(path)
		return nil
	})
}
