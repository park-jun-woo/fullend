//ff:func feature=gen-gogin type=util control=sequence topic=ddl
//ff:what SortTablesByFK — FK DAG 위상정렬 (Kahn's). 참조되는 테이블이 먼저.

package db

import "github.com/park-jun-woo/fullend/pkg/parser/ddl"

// SortTablesByFK returns table names in FK dependency order
// (referenced tables first). Returns error on FK cycle.
// Output is deterministic: alphabetical tie-breaker within the same depth.
func SortTablesByFK(tables []ddl.Table) ([]string, error) {
	inDegree, deps := buildFKGraph(tables)
	return topoSort(inDegree, deps)
}
