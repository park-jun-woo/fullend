//ff:func feature=symbol type=parser control=sequence topic=ddl
//ff:what DDL 한 줄을 파싱하여 테이블/컬럼/제약조건 정보를 tables에 반영한다

package validator

import "strings"

// parseDDLLine processes a single DDL line and updates tables accordingly.
// Returns the updated currentTable name.
func parseDDLLine(line string, currentTable string, tables map[string]DDLTable) string {
	line = strings.TrimSpace(line)
	upper := strings.ToUpper(line)

	if strings.HasPrefix(upper, "CREATE INDEX") || strings.HasPrefix(upper, "CREATE UNIQUE INDEX") {
		parseCreateIndex(line, tables)
		return currentTable
	}

	if strings.HasPrefix(upper, "CREATE TABLE") {
		return extractAndRegisterTable(line, tables)
	}

	if currentTable == "" {
		return currentTable
	}

	if strings.HasPrefix(line, ")") {
		return ""
	}

	if strings.HasPrefix(upper, "CONSTRAINT") || strings.HasPrefix(upper, "FOREIGN") {
		appendConstraintFK(line, currentTable, tables)
		return currentTable
	}

	if strings.HasPrefix(upper, "PRIMARY") {
		setPrimaryKey(line, currentTable, tables)
		return currentTable
	}

	if strings.HasPrefix(upper, "UNIQUE") {
		appendUniqueIndex(line, currentTable, tables)
		return currentTable
	}

	if line == "" {
		return currentTable
	}

	if strings.HasPrefix(upper, "CHECK") {
		applyCheckEnum(line, "", currentTable, tables)
		return currentTable
	}

	parts := strings.Fields(line)
	if len(parts) < 2 {
		return currentTable
	}

	colName := parts[0]
	colType := strings.ToUpper(parts[1])
	colType = strings.TrimSuffix(colType, ",")

	goType := pgTypeToGo(colType)
	t, ok := tables[currentTable]
	if !ok {
		return currentTable
	}
	t.Columns[colName] = goType
	t.ColumnOrder = append(t.ColumnOrder, colName)
	applyInlineConstraints(&t, upper, colName, parts)
	applyVarcharLen(&t, colName, colType)
	tables[currentTable] = t
	if strings.Contains(upper, "CHECK") {
		applyCheckEnum(line, colName, currentTable, tables)
	}

	return currentTable
}
