//ff:func feature=gen-gogin type=util control=iteration dimension=2 topic=ddl
//ff:what buildFKGraph — FK 관계를 위상정렬용 DAG 로 변환 (inDegree + 역방향 의존 목록)

package db

import "github.com/park-jun-woo/fullend/pkg/parser/ddl"

// buildFKGraph returns:
//   inDegree: tbl → count of tables it references (FK points outward)
//   deps    : refTable → list of tables that reference it (reverse for Kahn propagation)
func buildFKGraph(tables []ddl.Table) (map[string]int, map[string][]string) {
	names := map[string]bool{}
	for _, t := range tables {
		names[t.Name] = true
	}
	inDegree := map[string]int{}
	deps := map[string][]string{}
	for _, t := range tables {
		if _, ok := inDegree[t.Name]; !ok {
			inDegree[t.Name] = 0
		}
		seen := map[string]bool{}
		for _, fk := range t.ForeignKeys {
			if !names[fk.RefTable] || fk.RefTable == t.Name || seen[fk.RefTable] {
				continue
			}
			seen[fk.RefTable] = true
			deps[fk.RefTable] = append(deps[fk.RefTable], t.Name)
			inDegree[t.Name]++
		}
	}
	return inDegree, deps
}
