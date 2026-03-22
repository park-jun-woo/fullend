//ff:func feature=symbol type=loader control=iteration dimension=1 topic=sqlc
//ff:what 단일 SQL 파일에서 sqlc 메서드들을 파싱하여 ModelSymbol을 반환한다

package validator

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// parseSqlFileMethodsFromFile parses a single SQL query file and returns the model name and its methods.
func parseSqlFileMethodsFromFile(dir, fileName string) (string, ModelSymbol, error) {
	modelName := sqlFileToModel(fileName)
	ms := ModelSymbol{Methods: make(map[string]MethodInfo)}

	f, err := os.Open(filepath.Join(dir, fileName))
	if err != nil {
		return "", ms, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var currentMethod string
	var currentCardinality string
	var currentSQL strings.Builder

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "-- name:") && currentMethod != "" {
			currentSQL.WriteString(line + " ")
		}
		if !strings.HasPrefix(line, "-- name:") {
			continue
		}
		if currentMethod != "" {
			ms.Methods[currentMethod] = MethodInfo{
				Cardinality: currentCardinality,
				Params:      extractSqlcParams(currentSQL.String()),
			}
		}
		currentMethod, currentCardinality = parseSqlcNameLine(line, modelName)
		currentSQL.Reset()
	}
	if currentMethod != "" {
		params := extractSqlcParams(currentSQL.String())
		ms.Methods[currentMethod] = MethodInfo{
			Cardinality: currentCardinality,
			Params:      params,
		}
	}

	return modelName, ms, nil
}
