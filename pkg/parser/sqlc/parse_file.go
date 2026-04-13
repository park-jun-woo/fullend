//ff:func feature=sqlc-parse type=parser control=iteration dimension=1
//ff:what ParseFile — 단일 SQL 파일에서 sqlc 쿼리 항목 추출
package sqlc

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/fullend/pkg/diagnostic"
)

// ParseFile는 단일 `.sql` 파일에서 `-- name: Xxx :card` 항목을 추출한다.
func ParseFile(path string) ([]Query, []diagnostic.Diagnostic) {
	fileName := filepath.Base(path)
	model := sqlFileToModel(fileName)

	f, err := os.Open(path)
	if err != nil {
		return nil, []diagnostic.Diagnostic{{
			File:    path,
			Phase:   diagnostic.PhaseParse,
			Level:   diagnostic.LevelError,
			Message: "SQL 파일 열기 실패: " + err.Error(),
		}}
	}
	defer f.Close()

	var queries []Query
	scanner := bufio.NewScanner(f)

	var currentName, currentCard string
	var currentSQL strings.Builder

	flush := func() {
		if currentName == "" {
			return
		}
		queries = append(queries, Query{
			Model:       model,
			Name:        currentName,
			Cardinality: currentCard,
			Params:      extractParams(currentSQL.String()),
		})
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "-- name:") {
			if currentName != "" {
				currentSQL.WriteString(line)
				currentSQL.WriteByte(' ')
			}
			continue
		}
		flush()
		currentName, currentCard = parseNameLine(line, model)
		currentSQL.Reset()
	}
	flush()

	return queries, nil
}
