//ff:func feature=cli type=formatter control=sequence
//ff:what history 결과를 stdout에 출력 (캐시 fallback 포함)
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/clari/whyso/pkg/history"
)

// printHistoryStdout prints history results to stdout, falling back to cache if no new histories.
func printHistoryStdout(histories map[string]*history.FileHistory, absTarget, cwd, cacheDir, format string) {
	if len(histories) > 0 {
		printHistoryEntries(histories, format)
		return
	}
	// read from cache
	targetRel, _ := filepath.Rel(cwd, absTarget)
	cachedPath := filepath.Join(cacheDir, targetRel+"."+format)
	if cached, err := readHistoryYAML(cachedPath); err == nil {
		formatHistory(os.Stdout, cached, format)
		fmt.Println("---")
	}
}
