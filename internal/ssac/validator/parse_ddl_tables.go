//ff:func feature=symbol type=parser
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
			parts := strings.Fields(line)
			for i, p := range parts {
				pu := strings.ToUpper(p)
				if pu == "TABLE" && i+1 < len(parts) {
					currentTable = strings.Trim(parts[i+1], "( ")
					tables[currentTable] = DDLTable{Columns: make(map[string]string)}
					break
				}
			}
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
			if fk, ok := parseConstraintFK(line); ok {
				if t, exists := tables[currentTable]; exists {
					t.ForeignKeys = append(t.ForeignKeys, fk)
					tables[currentTable] = t
				}
			}
			continue
		}

		// PRIMARY KEY → PK 컬럼 추출
		if strings.HasPrefix(upper, "PRIMARY") {
			if t, ok := tables[currentTable]; ok {
				t.PrimaryKey = extractParenColumns(line)
				tables[currentTable] = t
			}
			continue
		}

		// UNIQUE 제약 (독립 라인) → unique index 추가
		if strings.HasPrefix(upper, "UNIQUE") {
			if t, ok := tables[currentTable]; ok {
				cols := extractParenColumns(line)
				if len(cols) > 0 {
					t.Indexes = append(t.Indexes, Index{Name: "unique_" + strings.Join(cols, "_"), Columns: cols, IsUnique: true})
					tables[currentTable] = t
				}
			}
			continue
		}

		// CHECK → skip
		if strings.HasPrefix(upper, "CHECK") || line == "" {
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
		if t, ok := tables[currentTable]; ok {
			t.Columns[colName] = goType
			t.ColumnOrder = append(t.ColumnOrder, colName)

			// 인라인 PRIMARY KEY
			if strings.Contains(upper, "PRIMARY KEY") {
				t.PrimaryKey = append(t.PrimaryKey, colName)
			}

			// 인라인 UNIQUE
			if strings.Contains(upper, "UNIQUE") && !strings.Contains(upper, "PRIMARY") {
				t.Indexes = append(t.Indexes, Index{Name: colName + "_unique", Columns: []string{colName}, IsUnique: true})
			}

			// 인라인 FK: column_name TYPE ... REFERENCES table(col)
			if fk, ok := parseInlineFK(colName, parts); ok {
				t.ForeignKeys = append(t.ForeignKeys, fk)
			}
			tables[currentTable] = t
		}
	}
}
