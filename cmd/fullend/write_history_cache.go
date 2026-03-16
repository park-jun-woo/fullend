//ff:func feature=cli type=util control=iteration dimension=1
//ff:what history 결과를 캐시 파일로 저장
package main

import (
	"os"
	"path/filepath"

	"github.com/clari/whyso/pkg/history"
)

// writeHistoryCache writes history entries to cache files, merging with existing data.
func writeHistoryCache(histories map[string]*history.FileHistory, cacheDir, format string) {
	for relPath, h := range histories {
		outPath := filepath.Join(cacheDir, relPath+"."+format)
		os.MkdirAll(filepath.Dir(outPath), 0755)
		// merge with existing
		if existing, err := readHistoryYAML(outPath); err == nil {
			h = history.Merge(existing, h)
		}
		f, err := os.Create(outPath)
		if err != nil {
			continue
		}
		formatHistory(f, h, format)
		f.Close()
	}
}
