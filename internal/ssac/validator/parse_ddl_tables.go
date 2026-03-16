//ff:func feature=symbol type=parser control=iteration dimension=1 topic=ddl
//ff:what CREATE TABLE 문에서 컬럼명, 타입, FK, 인덱스를 추출한다
package validator

import "strings"

// parseDDLTables는 CREATE TABLE 문에서 컬럼명, 타입, FK, 인덱스를 추출한다.
func parseDDLTables(content string, tables map[string]DDLTable) {
	lines := strings.Split(content, "\n")
	var currentTable string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		upper := strings.ToUpper(line)

		// CREATE INDEX idx_name ON tablename (col1, col2);
		if strings.HasPrefix(upper, "CREATE INDEX") || strings.HasPrefix(upper, "CREATE UNIQUE INDEX") {
			parseCreateIndex(line, tables)
			continue
		}

		// CREATE TABLE tablename (
		if strings.HasPrefix(upper, "CREATE TABLE") {
			currentTable = extractAndRegisterTable(line, tables)
			continue
		}

		if currentTable == "" {
			continue
		}

		// 테이블 정의 종료
		if strings.HasPrefix(line, ")") {
			currentTable = ""
			continue
		}

		// 독립 FOREIGN KEY: CONSTRAINT fk_name FOREIGN KEY (col) REFERENCES table(col)
		if strings.HasPrefix(upper, "CONSTRAINT") || strings.HasPrefix(upper, "FOREIGN") {
			appendConstraintFK(line, currentTable, tables)
			continue
		}

		// PRIMARY KEY → PK 컬럼 추출
		if strings.HasPrefix(upper, "PRIMARY") {
			setPrimaryKey(line, currentTable, tables)
			continue
		}

		// UNIQUE 제약 (독립 라인) → unique index 추가
		if strings.HasPrefix(upper, "UNIQUE") {
			appendUniqueIndex(line, currentTable, tables)
			continue
		}

		if line == "" {
			continue
		}

		// CHECK → parse enum values if present
		if strings.HasPrefix(upper, "CHECK") {
			applyCheckEnum(line, "", currentTable, tables)
			continue
		}

		// 컬럼 라인: column_name TYPE ...
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		colName := parts[0]
		colType := strings.ToUpper(parts[1])
		colType = strings.TrimSuffix(colType, ",")

		goType := pgTypeToGo(colType)
		t, ok := tables[currentTable]
		if !ok {
			continue
		}
		t.Columns[colName] = goType
		t.ColumnOrder = append(t.ColumnOrder, colName)
		applyInlineConstraints(&t, upper, colName, parts)
		applyVarcharLen(&t, colName, colType)
		if strings.Contains(upper, "CHECK") {
			applyCheckEnum(line, colName, currentTable, tables)
		}
		tables[currentTable] = t
	}
}
