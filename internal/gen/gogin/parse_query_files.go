//ff:func feature=gen-gogin type=parser control=iteration dimension=2
//ff:what parses query SQL files and extracts sqlc query annotations

package gogin

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// parseQueryFiles parses query SQL files from specsDir/db/queries/*.sql.
// Returns map[ModelName]map[MethodName]sqlcQuery.
func parseQueryFiles(specsDir string) map[string]map[string]sqlcQuery {
	result := make(map[string]map[string]sqlcQuery)

	queriesDir := filepath.Join(specsDir, "db", "queries")
	entries, err := os.ReadDir(queriesDir)
	if err != nil {
		return result
	}

	nameRe := regexp.MustCompile(`^--\s*name:\s*(\w+)\s+:(\w+)`)
	paramRe := regexp.MustCompile(`\$(\d+)`)
	insertColRe := regexp.MustCompile(`(?i)INSERT\s+INTO\s+\w+\s*\(([^)]+)\)`)
	updateSetRe := regexp.MustCompile(`(?i)SET\s+(.+?)(?:\s+WHERE|\s*;|\s*$)`)

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		baseName := strings.TrimSuffix(entry.Name(), ".sql")
		modelName := singularize(baseName)

		path := filepath.Join(queriesDir, entry.Name())
		queries := parseSingleQueryFile(path, modelName, nameRe, paramRe, insertColRe, updateSetRe)
		if len(queries) > 0 {
			result[modelName] = queries
		}
	}

	return result
}
