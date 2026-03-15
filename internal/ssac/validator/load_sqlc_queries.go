//ff:func feature=symbol type=loader
//ff:what queries/*.sql에서 모델과 메서드를 추출한다
package validator

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// loadSqlcQueries는 queries/*.sql에서 모델과 메서드를 추출한다.
// 파일명: users.sql → 모델 "User" (단수화 + PascalCase)
// 주석: -- name: FindByID :one → 메서드 "FindByID"
func (st *SymbolTable) loadSqlcQueries(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		modelName := sqlFileToModel(entry.Name())
		ms := ModelSymbol{Methods: make(map[string]MethodInfo)}

		f, err := os.Open(filepath.Join(dir, entry.Name()))
		if err != nil {
			return err
		}

		scanner := bufio.NewScanner(f)
		var currentMethod string
		var currentCardinality string
		var currentSQL strings.Builder

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			// -- name: FindByID :one  또는  -- name: CourseFindByID :one
			if strings.HasPrefix(line, "-- name:") {
				// 이전 메서드의 SQL 처리
				if currentMethod != "" {
					params := extractSqlcParams(currentSQL.String())
					ms.Methods[currentMethod] = MethodInfo{
						Cardinality: currentCardinality,
						Params:      params,
					}
				}
				// 새 메서드 시작
				parts := strings.Fields(line)
				if len(parts) >= 4 {
					currentMethod = stripModelPrefix(parts[2], modelName)
					currentCardinality = strings.TrimPrefix(parts[3], ":")
				} else if len(parts) >= 3 {
					currentMethod = stripModelPrefix(parts[2], modelName)
					currentCardinality = ""
				} else {
					currentMethod = ""
					currentCardinality = ""
				}
				currentSQL.Reset()
			} else if currentMethod != "" {
				currentSQL.WriteString(line + " ")
			}
		}
		// 마지막 메서드 처리
		if currentMethod != "" {
			params := extractSqlcParams(currentSQL.String())
			ms.Methods[currentMethod] = MethodInfo{
				Cardinality: currentCardinality,
				Params:      params,
			}
		}
		f.Close()

		if len(ms.Methods) > 0 {
			st.Models[modelName] = ms
		}
	}
	return nil
}
