//ff:func feature=gen-hurl type=parser
//ff:what Parses DDL CHECK constraints with IN (...) to extract enum values per column.
package hurl

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// parseDDLCheckEnums parses DDL files for CHECK constraints with IN (...) and returns
// a map of column_name -> allowed values.
// e.g. CHECK (plan_type IN ('free', 'pro', 'enterprise')) -> {"plan_type": ["free", "pro", "enterprise"]}
func parseDDLCheckEnums(specsDir string) map[string][]string {
	result := make(map[string][]string)
	dbDir := filepath.Join(specsDir, "db")
	entries, err := os.ReadDir(dbDir)
	if err != nil {
		return result
	}
	re := regexp.MustCompile(`CHECK\s*\(\s*(\w+)\s+IN\s*\(([^)]+)\)\s*\)`)
	valRe := regexp.MustCompile(`'([^']*)'`)
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dbDir, e.Name()))
		if err != nil {
			continue
		}
		for _, m := range re.FindAllStringSubmatch(string(data), -1) {
			col := m[1]
			var vals []string
			for _, vm := range valRe.FindAllStringSubmatch(m[2], -1) {
				vals = append(vals, vm[1])
			}
			if len(vals) > 0 {
				result[col] = vals
			}
		}
	}
	return result
}
