//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=string-convert
//ff:what map[string]bool의 키를 정렬하여 반환
package ssac

import "sort"

func sortedKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
