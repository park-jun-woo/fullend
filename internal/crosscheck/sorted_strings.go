//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=openapi-ddl
//ff:what map[string]bool의 키를 정렬하여 반환
package crosscheck

import "sort"

func sortedStrings(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
