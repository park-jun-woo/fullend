//ff:func feature=gen-gogin type=util control=iteration dimension=1 topic=ddl
//ff:what readDDLFileForTable — specsDir 에서 <tableName>.sql 파일 읽기 (case-insensitive)

package db

import (
	"os"
	"path/filepath"
	"strings"
)

// readDDLFileForTable finds <specsDir>/<name>.sql (case-insensitive) and returns its contents.
// Empty string if not found or read error.
func readDDLFileForTable(specsDir, tableName string) string {
	entries, err := os.ReadDir(specsDir)
	if err != nil {
		return ""
	}
	lower := strings.ToLower(tableName)
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		base := strings.TrimSuffix(strings.ToLower(e.Name()), ".sql")
		if base != lower {
			continue
		}
		data, err := os.ReadFile(filepath.Join(specsDir, e.Name()))
		if err == nil {
			return string(data)
		}
	}
	return ""
}
