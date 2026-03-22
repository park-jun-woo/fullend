//ff:func feature=gen-gogin type=parser control=iteration dimension=2
//ff:what parses CREATE TABLE statements from DDL files

package gogin

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// parseDDLFiles parses CREATE TABLE statements from specsDir/db/*.sql.
func parseDDLFiles(specsDir string) map[string]*ddlTable {
	tables := make(map[string]*ddlTable)

	dbDir := filepath.Join(specsDir, "db")
	entries, err := os.ReadDir(dbDir)
	if err != nil {
		return tables
	}

	createRe := regexp.MustCompile(`(?i)CREATE\s+TABLE\s+(\w+)\s*\(`)
	colRe := regexp.MustCompile(`^\s+(\w+)\s+(BIGSERIAL|BIGINT|INT|INTEGER|VARCHAR\(\d+\)|TEXT|BOOLEAN|BOOL|TIMESTAMPTZ|TIMESTAMP|JSONB|JSON)`)
	fkRe := regexp.MustCompile(`REFERENCES\s+(\w+)\s*\(`)

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		path := filepath.Join(dbDir, entry.Name())
		table := parseSingleDDLFile(path, createRe, colRe, fkRe)
		if table != nil {
			tables[table.ModelName] = table
		}
	}

	return tables
}
