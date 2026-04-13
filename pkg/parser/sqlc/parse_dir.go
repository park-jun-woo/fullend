//ff:func feature=sqlc-parse type=parser control=iteration dimension=1
//ff:what ParseDir — queries 디렉토리의 *.sql 파일을 순회하며 쿼리 추출
package sqlc

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/fullend/pkg/diagnostic"
)

// ParseDir는 주어진 디렉토리의 *.sql 파일을 순회하며 sqlc 쿼리를 추출한다.
// 디렉토리가 없으면 nil, nil 을 반환한다.
func ParseDir(dir string) ([]Query, []diagnostic.Diagnostic) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, []diagnostic.Diagnostic{{
			File:    dir,
			Phase:   diagnostic.PhaseParse,
			Level:   diagnostic.LevelError,
			Message: "queries 디렉토리 읽기 실패: " + err.Error(),
		}}
	}

	var queries []Query
	var diags []diagnostic.Diagnostic
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		fileQueries, fileDiags := ParseFile(path)
		diags = append(diags, fileDiags...)
		queries = append(queries, fileQueries...)
	}
	return queries, diags
}
