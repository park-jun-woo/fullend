//ff:func feature=gen-gogin type=util control=iteration dimension=1 topic=ddl
//ff:what decrementDependents — 처리된 테이블에 의존하는 테이블들의 inDegree 를 하나 감소

package db

func decrementDependents(name string, inDegree map[string]int, deps map[string][]string) {
	for _, dependent := range deps[name] {
		if inDegree[dependent] > 0 {
			inDegree[dependent]--
		}
	}
}
