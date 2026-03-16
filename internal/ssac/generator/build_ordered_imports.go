//ff:func feature=ssac-gen type=util control=iteration dimension=1
//ff:what 표준 라이브러리 우선 순서로 import 슬라이스를 생성
package generator

import "sort"

func buildOrderedImports(seen map[string]bool) []string {
	var imports []string
	order := []string{"database/sql", "encoding/json", "net/http", "strconv", "time"}
	for _, imp := range order {
		if seen[imp] {
			imports = append(imports, imp)
			delete(seen, imp)
		}
	}
	var dynamic []string
	for imp := range seen {
		dynamic = append(dynamic, imp)
	}
	sort.Strings(dynamic)
	return append(imports, dynamic...)
}
