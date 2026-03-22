//ff:func feature=gen-gogin type=parser control=iteration dimension=1
//ff:what 쿼리 SQL 파일 하나를 파싱하여 sqlcQuery 맵을 반환한다

package gogin

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

// parseSingleQueryFile parses a single query SQL file and returns a map of method name to sqlcQuery.
func parseSingleQueryFile(path, modelName string, nameRe, paramRe, insertColRe, updateSetRe *regexp.Regexp) map[string]sqlcQuery {
	queries := make(map[string]sqlcQuery)

	f, err := os.Open(path)
	if err != nil {
		return queries
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var currentQuery *sqlcQuery
	var sqlBuf strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		matches := nameRe.FindStringSubmatch(line)
		if matches == nil && currentQuery != nil {
			sqlBuf.WriteString(line)
			sqlBuf.WriteString("\n")
		}
		if matches == nil {
			continue
		}
		if currentQuery != nil {
			finishQuery(currentQuery, sqlBuf.String(), paramRe, insertColRe, updateSetRe)
			queries[currentQuery.Name] = *currentQuery
		}
		currentQuery = &sqlcQuery{
			Name:        stripModelPrefix(matches[1], modelName),
			Cardinality: matches[2],
		}
		sqlBuf.Reset()
	}
	if currentQuery != nil {
		finishQuery(currentQuery, sqlBuf.String(), paramRe, insertColRe, updateSetRe)
		queries[currentQuery.Name] = *currentQuery
	}

	return queries
}
