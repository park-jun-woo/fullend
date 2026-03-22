//ff:func feature=symbol type=loader control=iteration dimension=2 topic=sqlc
//ff:what queries/*.sql에서 모델과 메서드를 추출한다
package validator

import (
	"os"
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

		modelName, ms, err := parseSqlFileMethodsFromFile(dir, entry.Name())
		if err != nil {
			return err
		}
		if len(ms.Methods) > 0 {
			st.Models[modelName] = ms
		}
	}
	return nil
}

// parseSqlcNameLine은 "-- name: FindByID :one" 형식의 줄에서 메서드명과 cardinality를 추출한다.
func parseSqlcNameLine(line, modelName string) (method, cardinality string) {
	parts := strings.Fields(line)
	if len(parts) >= 4 {
		return stripModelPrefix(parts[2], modelName), strings.TrimPrefix(parts[3], ":")
	}
	if len(parts) >= 3 {
		return stripModelPrefix(parts[2], modelName), ""
	}
	return "", ""
}
