//ff:func feature=rule type=util control=sequence
//ff:what ensureForbiddenSets — Ground에 금지 참조 목록이 없으면 등록
package ssac

import "github.com/park-jun-woo/fullend/pkg/rule"

func ensureForbiddenSets(g *rule.Ground) {
	if _, ok := g.Lookup["ssac.reservedSource"]; !ok {
		g.Lookup["ssac.reservedSource"] = rule.StringSet{
			"request": true, "currentUser": true, "config": true, "query": true, "message": true,
		}
	}
	if _, ok := g.Lookup["ssac.configPrefix"]; !ok {
		g.Lookup["ssac.configPrefix"] = rule.StringSet{"config": true}
	}
}
