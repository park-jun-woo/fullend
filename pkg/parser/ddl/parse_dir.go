//ff:func feature=orchestrator type=parser control=iteration dimension=1
//ff:what 디렉토리 내 모든 .sql 파일을 pg_query_go로 파싱
package ddl

import (
	"os"
	"path/filepath"
	"strings"

	pg_query "github.com/pganalyze/pg_query_go/v5"

	"github.com/park-jun-woo/fullend/pkg/diagnostic"
)

// ParseDir parses all .sql files in the given directory using pg_query_go.
func ParseDir(dir string) ([]*pg_query.ParseResult, []diagnostic.Diagnostic) {
	var results []*pg_query.ParseResult
	var diags []diagnostic.Diagnostic

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, []diagnostic.Diagnostic{{
			File:    dir,
			Line:    0,
			Phase:   diagnostic.PhaseParse,
			Level:   diagnostic.LevelError,
			Message: "cannot read DDL directory: " + err.Error(),
		}}
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		filePath := filepath.Join(dir, e.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			diags = append(diags, diagnostic.Diagnostic{
				File:    filePath,
				Line:    0,
				Phase:   diagnostic.PhaseParse,
				Level:   diagnostic.LevelError,
				Message: "cannot read SQL file: " + err.Error(),
			})
			continue
		}
		result, err := pg_query.Parse(string(data))
		if err != nil {
			diags = append(diags, diagnostic.Diagnostic{
				File:    filePath,
				Line:    0,
				Phase:   diagnostic.PhaseParse,
				Level:   diagnostic.LevelError,
				Message: "SQL parse error: " + err.Error(),
			})
			continue
		}
		results = append(results, result)
	}
	return results, diags
}
