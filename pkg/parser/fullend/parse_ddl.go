//ff:func feature=orchestrator type=parser control=iteration dimension=1
//ff:what 디렉토리 내 모든 .sql 파일을 pg_query_go로 파싱
package fullend

import (
	"os"
	"path/filepath"
	"strings"

	pg_query "github.com/pganalyze/pg_query_go/v5"
)

// parseDDLDir parses all .sql files in the given directory using pg_query_go.
func parseDDLDir(dir string) []*pg_query.ParseResult {
	var results []*pg_query.ParseResult
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		result, err := pg_query.Parse(string(data))
		if err != nil {
			continue
		}
		results = append(results, result)
	}
	return results
}
