//ff:func feature=cli type=util control=sequence
//ff:what history 대상 필터 함수 생성
package main

import (
	"os"
	"path/filepath"
	"strings"
)

// buildHistoryFilter returns a filter function that matches paths against the target.
func buildHistoryFilter(absTarget, cwd string, targetInfo os.FileInfo, all bool) func(string) bool {
	return func(relPath string) bool {
		if targetInfo.IsDir() {
			if !all {
				return false
			}
			targetRel, err := filepath.Rel(cwd, absTarget)
			if err != nil {
				return false
			}
			if targetRel == "." {
				return true
			}
			return strings.HasPrefix(relPath, targetRel+"/") || relPath == targetRel
		}
		targetRel, err := filepath.Rel(cwd, absTarget)
		if err != nil {
			return false
		}
		return relPath == targetRel
	}
}
