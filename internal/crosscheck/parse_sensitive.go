//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what DDL 디렉토리에서 @sensitive/@nosensitive 태그를 파싱
package crosscheck

import (
	"os"
	"path/filepath"
	"strings"
)

// ParseSensitive parses DDL .sql files in dbDir for @sensitive and @nosensitive tags.
func ParseSensitive(dbDir string) (sensitive, nosensitive map[string]map[string]bool, err error) {
	sensitive = make(map[string]map[string]bool)
	nosensitive = make(map[string]map[string]bool)

	entries, err := os.ReadDir(dbDir)
	if err != nil {
		if os.IsNotExist(err) {
			return sensitive, nosensitive, nil
		}
		return nil, nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dbDir, entry.Name()))
		if err != nil {
			return nil, nil, err
		}
		parseSensitiveSQL(string(data), sensitive, nosensitive)
	}

	return sensitive, nosensitive, nil
}
