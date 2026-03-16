//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=args-inputs
//ff:what paramOrder 순서대로 입력 키를 배치하고 나머지를 알파벳순 추가
package generator

import "sort"

func orderByParamOrder(inputs map[string]string, paramOrder []string) []string {
	var keys []string
	used := make(map[string]bool)
	for _, p := range paramOrder {
		if _, ok := inputs[p]; ok {
			keys = append(keys, p)
			used[p] = true
		}
	}
	var extra []string
	for k := range inputs {
		if !used[k] {
			extra = append(extra, k)
		}
	}
	sort.Strings(extra)
	return append(keys, extra...)
}
