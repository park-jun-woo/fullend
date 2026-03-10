package crosscheck

import (
	"os"
	"path/filepath"
	"strings"
)

// ArchivedInfo holds @archived tags parsed from DDL files.
type ArchivedInfo struct {
	Tables  map[string]bool            // "legacy_notifications" → true
	Columns map[string]map[string]bool // "courses" → {"old_category": true}
}

// ParseArchived parses DDL .sql files in dbDir for @archived tags.
// Table-level: "-- @archived" on the line before CREATE TABLE.
// Column-level: "-- @archived" at the end of a column definition line.
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

func parseArchivedSQL(content string, info *ArchivedInfo) {
	lines := strings.Split(content, "\n")
	prevLineArchived := false
	var currentTable string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		upper := strings.ToUpper(trimmed)

		// Detect standalone "-- @archived" comment.
		if strings.HasPrefix(trimmed, "--") && strings.Contains(trimmed, "@archived") {
			// Check if this is a standalone line (not inline on a column).
			// Standalone: the line is ONLY a comment.
			prevLineArchived = true
			continue
		}

		// CREATE TABLE
		if strings.HasPrefix(upper, "CREATE TABLE") {
			parts := strings.Fields(trimmed)
			for i, p := range parts {
				if strings.ToUpper(p) == "TABLE" && i+1 < len(parts) {
					currentTable = strings.Trim(parts[i+1], "( ")
					break
				}
			}
			if prevLineArchived && currentTable != "" {
				info.Tables[currentTable] = true
			}
			prevLineArchived = false
			continue
		}

		prevLineArchived = false

		if currentTable == "" {
			continue
		}

		// End of table definition.
		if strings.HasPrefix(trimmed, ")") {
			currentTable = ""
			continue
		}

		// Skip constraints, primary key, etc.
		if strings.HasPrefix(upper, "PRIMARY") || strings.HasPrefix(upper, "UNIQUE") ||
			strings.HasPrefix(upper, "CHECK") || strings.HasPrefix(upper, "CONSTRAINT") ||
			strings.HasPrefix(upper, "FOREIGN") || trimmed == "" {
			continue
		}

		// Column line with inline -- @archived.
		if strings.Contains(line, "-- @archived") {
			colParts := strings.Fields(trimmed)
			if len(colParts) >= 2 {
				colName := colParts[0]
				if info.Columns[currentTable] == nil {
					info.Columns[currentTable] = make(map[string]bool)
				}
				info.Columns[currentTable][colName] = true
			}
		}
	}
}
