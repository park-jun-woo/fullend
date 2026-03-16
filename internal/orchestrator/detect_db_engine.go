//ff:func feature=orchestrator type=util control=iteration
//ff:what detectDBEngine inspects DDL files to determine the database engine.

package orchestrator

import (
	"os"
	"path/filepath"
	"strings"
)

// detectDBEngine inspects DDL files to determine the database engine.
func detectDBEngine(specsDir string) string {
	dbDir := filepath.Join(specsDir, "db")
	entries, err := os.ReadDir(dbDir)
	if err != nil {
		return "postgresql"
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dbDir, entry.Name()))
		if err != nil {
			continue
		}
		content := strings.ToUpper(string(data))
		// MySQL indicators.
		if strings.Contains(content, "AUTO_INCREMENT") ||
			strings.Contains(content, "ENGINE=INNODB") ||
			strings.Contains(content, "ENGINE = INNODB") {
			return "mysql"
		}
	}

	return "postgresql"
}
