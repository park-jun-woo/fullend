//ff:func feature=manifest type=parser control=sequence
//ff:what parseColumnDef — 컬럼 정의 라인에서 이름, 타입, 인라인 제약 추출
package ddl

import "strings"

func parseColumnDef(line, upper string, t *Table) {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return
	}
	colName := parts[0]
	colType := strings.ToUpper(parts[1])
	colType = strings.TrimSuffix(colType, ",")

	t.Columns[colName] = pgTypeToGo(colType)
	t.ColumnOrder = append(t.ColumnOrder, colName)
	applyInlineConstraints(t, upper, colName, parts)
	applyVarcharLen(t, colName, colType)
	if strings.Contains(upper, "CHECK") {
		applyCheckEnum(line, colName, t)
	}
}
