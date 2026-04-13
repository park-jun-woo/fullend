//ff:func feature=gen-gogin type=parser control=iteration dimension=1
//ff:what DDL 파일 하나를 파싱하여 테이블 컬럼 정보를 추출한다

package gogin

import (
	"os"
	"regexp"
	"strings"
)

// parseSingleDDLFile parses a single DDL file and returns the ddlTable, or nil if not applicable.
func parseSingleDDLFile(path string, createRe, colRe, fkRe *regexp.Regexp) *ddlTable {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	content := string(data)
	tableMatch := createRe.FindStringSubmatch(content)
	if tableMatch == nil {
		return nil
	}

	tableName := tableMatch[1]
	modelName := singularize(tableName)

	table := &ddlTable{
		TableName: tableName,
		ModelName: modelName,
	}

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		col := parseDDLColumn(line, colRe, fkRe)
		if col != nil {
			table.Columns = append(table.Columns, *col)
		}
	}

	return table
}
