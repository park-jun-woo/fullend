//ff:func feature=gen-gogin type=parser control=sequence
//ff:what DDL 라인에서 컬럼 정보를 추출하여 ddlColumn 반환
package gogin

import (
	"regexp"
	"strings"
)

// parseDDLColumn extracts a ddlColumn from a DDL line, or returns nil if not a column line.
func parseDDLColumn(line string, colRe, fkRe *regexp.Regexp) *ddlColumn {
	colMatch := colRe.FindStringSubmatch(line)
	if colMatch == nil {
		return nil
	}
	colName := colMatch[1]
	sqlType := strings.ToUpper(colMatch[2])
	fkTable := ""
	if fkMatch := fkRe.FindStringSubmatch(line); fkMatch != nil {
		fkTable = fkMatch[1]
	}
	upperLine := strings.ToUpper(line)
	notNull := strings.Contains(upperLine, "NOT NULL") || strings.Contains(upperLine, "PRIMARY KEY")
	sensitive := strings.Contains(line, "@sensitive")
	return &ddlColumn{
		Name:      colName,
		GoName:    snakeToGo(colName),
		GoType:    sqlTypeToGo(sqlType),
		FKTable:   fkTable,
		NotNull:   notNull,
		Sensitive: sensitive,
	}
}
