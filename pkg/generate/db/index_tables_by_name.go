//ff:func feature=gen-gogin type=util control=iteration dimension=1 topic=ddl
//ff:what indexTablesByName — []ddl.Table 을 이름→*Table 맵으로

package db

import "github.com/park-jun-woo/fullend/pkg/parser/ddl"

func indexTablesByName(tables []ddl.Table) map[string]*ddl.Table {
	m := make(map[string]*ddl.Table, len(tables))
	for i := range tables {
		m[tables[i].Name] = &tables[i]
	}
	return m
}
