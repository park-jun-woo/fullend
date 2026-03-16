//ff:func feature=gen-hurl type=parser control=iteration dimension=2
//ff:what DDL 파일에서 FK 관계를 파싱한다
package hurl

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// parseDDLFiles extracts table names and FK references from SQL DDL files.
// This is a minimal subset of gogin's full DDL parser — only FK deps are needed here.
func parseDDLFiles(specsDir string) map[string]*ddlFK {
	tables := make(map[string]*ddlFK)

	dbDir := filepath.Join(specsDir, "db")
	entries, err := os.ReadDir(dbDir)
	if err != nil {
		return tables
	}

	createRe := regexp.MustCompile(`(?i)CREATE\s+TABLE\s+(\w+)\s*\(`)
	fkRe := regexp.MustCompile(`REFERENCES\s+(\w+)\s*\(`)

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dbDir, entry.Name()))
		if err != nil {
			continue
		}

		content := string(data)
		tableMatch := createRe.FindStringSubmatch(content)
		if tableMatch == nil {
			continue
		}

		tableName := tableMatch[1]
		t := &ddlFK{TableName: tableName}

		for _, line := range strings.Split(content, "\n") {
			if fkMatch := fkRe.FindStringSubmatch(line); fkMatch != nil {
				t.FKTables = append(t.FKTables, fkMatch[1])
			}
		}

		tables[tableName] = t
	}

	return tables
}
