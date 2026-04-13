//ff:func feature=gen-gogin type=util control=iteration dimension=1 topic=ddl
//ff:what pickZeroDegree — inDegree 맵에서 degree=0 인 테이블 이름 수집 (위상정렬 step)

package db

func pickZeroDegree(inDegree map[string]int) []string {
	var ready []string
	for name, d := range inDegree {
		if d == 0 {
			ready = append(ready, name)
		}
	}
	return ready
}
