//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=ssac-ddl
//ff:what foreignKeyMap — ForeignKey 슬라이스를 컬럼명→FK 맵으로

package crosscheck

import "github.com/park-jun-woo/fullend/pkg/parser/ddl"

func foreignKeyMap(fks []ddl.ForeignKey) map[string]ddl.ForeignKey {
	m := map[string]ddl.ForeignKey{}
	for _, fk := range fks {
		m[fk.Column] = fk
	}
	return m
}
