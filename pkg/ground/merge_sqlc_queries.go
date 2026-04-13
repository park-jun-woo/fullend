//ff:func feature=rule type=loader control=iteration dimension=1
//ff:what mergeSqlcQueries — sqlc.Query 배열로 모델 메서드에 cardinality + params 주입
package ground

import (
	"github.com/park-jun-woo/fullend/pkg/parser/sqlc"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func mergeSqlcQueries(g *rule.Ground, queries []sqlc.Query) {
	for _, q := range queries {
		info := ensureModel(g, q.Model)
		info.Methods[q.Name] = rule.MethodInfo{
			Cardinality: q.Cardinality,
			Params:      q.Params,
		}
		g.Models[q.Model] = info
	}
}
