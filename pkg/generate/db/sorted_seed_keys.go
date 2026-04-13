//ff:func feature=gen-gogin type=util control=iteration dimension=1 topic=ddl
//ff:what sortedSeedKeys — seed 키 set 을 정렬된 슬라이스로

package db

import "sort"

func sortedSeedKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
