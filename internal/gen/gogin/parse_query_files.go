//ff:func feature=gen-gogin type=parser control=iteration
//ff:what parses query SQL files and extracts sqlc query annotations

package gogin

import (
	"bufio"
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

		// Derive model name from filename (e.g. "course.sql" -> "Course").
		baseName := strings.TrimSuffix(entry.Name(), ".sql")
		modelName := singularize(baseName)

		if result[modelName] == nil {
			result[modelName] = make(map[string]sqlcQuery)
		}

		path := filepath.Join(queriesDir, entry.Name())
		f, err := os.Open(path)
		if err != nil {
			continue
		}

		scanner := bufio.NewScanner(f)
		var currentQuery *sqlcQuery
		var sqlBuf strings.Builder

		for scanner.Scan() {
			line := scanner.Text()

			if matches := nameRe.FindStringSubmatch(line); matches != nil {
				// Save previous query.
				if currentQuery != nil {
					finishQuery(currentQuery, sqlBuf.String(), paramRe, insertColRe, updateSetRe)
					result[modelName][currentQuery.Name] = *currentQuery
				}
				currentQuery = &sqlcQuery{
					Name:        stripModelPrefix(matches[1], modelName),
					Cardinality: matches[2],
				}
				sqlBuf.Reset()
			} else if currentQuery != nil {
				sqlBuf.WriteString(line)
				sqlBuf.WriteString("\n")
			}
		}
		// Save last query in file.
		if currentQuery != nil {
			finishQuery(currentQuery, sqlBuf.String(), paramRe, insertColRe, updateSetRe)
			result[modelName][currentQuery.Name] = *currentQuery
		}
		f.Close()
	}

	return result
}
