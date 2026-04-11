//ff:func feature=manifest type=parser control=iteration dimension=1
//ff:what ParseTables — db/ 디렉토리의 .sql 파일에서 Table 목록 추출
package ddl

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/fullend/pkg/diagnostic"
)

// ParseTables reads all .sql files in dir and returns parsed tables.
func ParseTables(dir string) ([]Table, []diagnostic.Diagnostic) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, []diagnostic.Diagnostic{{Message: "cannot read DDL dir: " + err.Error()}}
	}
	tables := make(map[string]*Table)
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		parseDDLContent(string(data), tables)
	}
	var result []Table
	for _, t := range tables {
		result = append(result, *t)
	}
	return result, nil
}
