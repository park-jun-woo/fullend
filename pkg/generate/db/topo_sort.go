//ff:func feature=gen-gogin type=util control=iteration dimension=1 topic=ddl
//ff:what topoSort — Kahn's algorithm 위상정렬 (결정적: depth 내 알파벳순)

package db

import "fmt"

func topoSort(inDegree map[string]int, deps map[string][]string) ([]string, error) {
	remaining := len(inDegree)
	var result []string
	for remaining > 0 {
		ready := pickZeroDegree(inDegree)
		if len(ready) == 0 {
			return nil, fmt.Errorf("FK cycle detected among tables: %v", pendingTables(inDegree))
		}
		result, remaining = applyZeroDegreeBatch(ready, result, inDegree, deps, remaining)
	}
	return result, nil
}
