//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what DDL 디렉토리에서 @archived 태그를 파싱
package crosscheck

import (
	"os"
	"path/filepath"
	"strings"
)

// ParseArchived parses DDL .sql files in dbDir for @archived tags.
func ParseArchived(dbDir string) (*ArchivedInfo, error) {
	info := &ArchivedInfo{
		Tables:  make(map[string]bool),
		Columns: make(map[string]map[string]bool),
	}

	entries, err := os.ReadDir(dbDir)
	if err != nil {
		if os.IsNotExist(err) {
			return info, nil
		}
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dbDir, entry.Name()))
		if err != nil {
			return nil, err
		}
		parseArchivedSQL(string(data), info)
	}

	return info, nil
}
