//ff:func feature=gen-gogin type=util control=iteration dimension=1 topic=ddl
//ff:what pendingTables — 아직 처리 안 된 (degree>0) 테이블 이름. 사이클 에러 메시지용

package db

import "sort"

func pendingTables(inDegree map[string]int) []string {
	var out []string
	for name, d := range inDegree {
		if d > 0 {
			out = append(out, name)
		}
	}
	sort.Strings(out)
	return out
}
