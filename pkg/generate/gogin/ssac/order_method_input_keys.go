//ff:func feature=ssac-gen type=util control=selection topic=args-inputs
//ff:what 메서드 파라미터 순서에 따라 입력 키를 정렬
package ssac

import "sort"

func orderMethodInputKeys(inputs map[string]string, paramOrder []string) []string {
	switch {
	case len(paramOrder) > 0:
		return orderInputsByParams(inputs, paramOrder)
	default:
		keys := make([]string, 0, len(inputs))
		for k := range inputs {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		return keys
	}
}
