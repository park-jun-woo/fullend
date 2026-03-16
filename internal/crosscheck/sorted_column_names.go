//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what DDL 컬럼 맵에서 정렬된 컬럼명 목록 반환
package crosscheck

import "sort"

// sortedColumnNames returns sorted column names from a DDL column map.
func sortedColumnNames(columns map[string]string) []string {
	keys := make([]string, 0, len(columns))
	for col := range columns {
		keys = append(keys, col)
	}
	sort.Strings(keys)
	return keys
}
