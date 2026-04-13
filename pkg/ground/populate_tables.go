//ff:func feature=rule type=loader control=iteration dimension=1
//ff:what populateTables — DDLTables 를 g.Tables 로 복사 (Columns + ColumnOrder 보존)
package ground

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func populateTables(g *rule.Ground, fs *fullend.Fullstack) {
	if g.Tables == nil {
		g.Tables = make(map[string]rule.TableInfo)
	}
	for _, t := range fs.DDLTables {
		columns := make(map[string]string, len(t.Columns))
		for k, v := range t.Columns {
			columns[k] = v
		}
		order := make([]string, len(t.ColumnOrder))
		copy(order, t.ColumnOrder)
		g.Tables[t.Name] = rule.TableInfo{
			Name:        t.Name,
			Columns:     columns,
			ColumnOrder: order,
		}
	}
}
