//ff:func feature=symbol type=util control=sequence
//ff:what 컬럼 라인의 인라인 제약(PK, UNIQUE, FK)을 DDLTable에 반영한다
package validator

import "strings"

func applyInlineConstraints(t *DDLTable, upper, colName string, parts []string) {
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
}
