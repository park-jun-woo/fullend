//ff:func feature=gen-gogin type=util control=iteration dimension=1 topic=ddl
//ff:what applyZeroDegreeBatch — ready 테이블들을 결과에 추가하고 의존 감소 처리

package db

import "sort"

// applyZeroDegreeBatch visits ready tables, appends to result, marks visited (-1),
// and propagates inDegree decrements. Returns updated result + remaining count.
func applyZeroDegreeBatch(ready []string, result []string, inDegree map[string]int, deps map[string][]string, remaining int) ([]string, int) {
	sort.Strings(ready)
	for _, name := range ready {
		result = append(result, name)
		inDegree[name] = -1
		remaining--
		decrementDependents(name, inDegree, deps)
	}
	return result, remaining
}
