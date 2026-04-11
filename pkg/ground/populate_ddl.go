//ff:func feature=rule type=loader control=iteration dimension=1
//ff:what populateDDL — DDL Table에서 테이블명, 컬럼, FK, 인덱스, CHECK 추출
package ground

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func populateDDL(g *rule.Ground, fs *fullend.Fullstack) {
	tables := make(rule.StringSet)
	for _, t := range fs.DDLTables {
		tables[t.Name] = true
		cols := make(rule.StringSet, len(t.Columns))
		for col := range t.Columns {
			cols[col] = true
		}
		g.Lookup["DDL.column."+t.Name] = cols
		populateDDLIndexes(g, t)
		populateDDLCheck(g, t)
		populateDDLVarchar(g, t)
	}
	g.Lookup["DDL.table"] = tables
}
