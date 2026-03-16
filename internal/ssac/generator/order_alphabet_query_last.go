//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=args-inputs
//ff:what 입력 키를 알파벳순으로 정렬하되 query를 마지막에 배치
package generator

import "sort"

func orderAlphabetQueryLast(inputs map[string]string) []string {
	var keys []string
	var queryKey string
	for k := range inputs {
		if inputs[k] == "query" {
			queryKey = k
		} else {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	if queryKey != "" {
		keys = append(keys, queryKey)
	}
	return keys
}
