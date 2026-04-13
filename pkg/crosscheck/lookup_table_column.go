//ff:func feature=crosscheck type=util control=sequence topic=config-check
//ff:what lookupTableColumn — Ground.Tables[table].Columns[col] 조회 (flat 가드)

package crosscheck

import "github.com/park-jun-woo/fullend/pkg/rule"

func lookupTableColumn(table, col string, g *rule.Ground) (string, bool) {
	t, hit := g.Tables[table]
	if !hit {
		return "", false
	}
	gt, has := t.Columns[col]
	return gt, has
}
