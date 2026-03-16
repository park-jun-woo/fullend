//ff:func feature=gen-gogin type=parser control=iteration
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
	// Match column definitions: name TYPE(...) constraints
	// Stop at lines starting with constraints or indexes.
	colRe := regexp.MustCompile(`^\s+(\w+)\s+(BIGSERIAL|BIGINT|INT|INTEGER|VARCHAR\(\d+\)|TEXT|BOOLEAN|BOOL|TIMESTAMPTZ|TIMESTAMP|JSONB|JSON)`)
	fkRe := regexp.MustCompile(`REFERENCES\s+(\w+)\s*\(`)

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		path := filepath.Join(dbDir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		content := string(data)
		tableMatch := createRe.FindStringSubmatch(content)
		if tableMatch == nil {
			continue
		}

		tableName := tableMatch[1]
		modelName := singularize(tableName)

		table := &ddlTable{
			TableName: tableName,
			ModelName: modelName,
		}

		lines := strings.Split(content, "\n")
		for _, line := range lines {
			colMatch := colRe.FindStringSubmatch(line)
			if colMatch == nil {
				continue
			}
			colName := colMatch[1]
			sqlType := strings.ToUpper(colMatch[2])

			fkTable := ""
			if fkMatch := fkRe.FindStringSubmatch(line); fkMatch != nil {
				fkTable = fkMatch[1]
			}

			upperLine := strings.ToUpper(line)
			notNull := strings.Contains(upperLine, "NOT NULL") || strings.Contains(upperLine, "PRIMARY KEY")
			sensitive := strings.Contains(line, "@sensitive")

			table.Columns = append(table.Columns, ddlColumn{
				Name:      colName,
				GoName:    snakeToGo(colName),
				GoType:    sqlTypeToGo(sqlType),
				FKTable:   fkTable,
				NotNull:   notNull,
				Sensitive: sensitive,
			})
		}

		tables[modelName] = table
	}

	return tables
}
