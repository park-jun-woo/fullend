//ff:func feature=gen-gogin type=util control=selection topic=ddl
//ff:what seedValueFor — 테이블 컬럼 하나에 대한 SQL 리터럴 값 (CHECK/타입 기반)

package db

import (
	"strconv"
	"strings"

	"github.com/park-jun-woo/fullend/pkg/parser/ddl"
)

// seedValueFor returns a SQL literal for one column's auto-seed value.
func seedValueFor(t *ddl.Table, col string, idVal int64) string {
	if col == "id" {
		return strconv.FormatInt(idVal, 10)
	}
	if enums, has := t.CheckEnums[col]; has && len(enums) > 0 {
		return "'" + enums[0] + "'"
	}
	goType := t.Columns[col]
	switch {
	case strings.HasPrefix(goType, "int") || goType == "float64":
		return "0"
	case goType == "bool":
		return "false"
	default:
		return "'" + nobodyPlaceholderFor(t.Name, col) + "'"
	}
}
