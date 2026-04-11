//ff:func feature=crosscheck type=rule control=sequence
//ff:what checkResultDDLTable — @result 타입에 대응하는 DDL 테이블 존재 여부 검증 (X-12)
package crosscheck

import (
	"strings"

	"github.com/jinzhu/inflection"

	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkResultDDLTable(g *rule.Ground, funcName, resultType string) []CrossError {
	table := strings.ToLower(inflection.Plural(resultType))
	if _, ok := g.Lookup["DDL.column."+table]; ok {
		return nil
	}
	if g.Lookup["DDL.table"][table] {
		return nil
	}
	return []CrossError{{Rule: "X-12", Context: funcName, Level: "WARNING",
		Message: "@result type " + resultType + " has no matching DDL table"}}
}
