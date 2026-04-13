//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=ssac-ddl
//ff:what tableMapByName — ddl.Table 슬라이스를 이름→*Table 맵으로

package crosscheck

import "github.com/park-jun-woo/fullend/pkg/parser/ddl"

func tableMapByName(ts []ddl.Table) map[string]*ddl.Table {
	m := map[string]*ddl.Table{}
	for i := range ts {
		m[ts[i].Name] = &ts[i]
	}
	return m
}
