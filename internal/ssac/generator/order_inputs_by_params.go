//ff:func feature=ssac-gen type=util control=iteration dimension=1
//ff:what sqlc/interface 파라미터 순서대로 입력 키를 배치
package generator

func orderInputsByParams(inputs map[string]string, paramOrder []string) []string {
	var keys []string
	used := make(map[string]bool)
	for _, p := range paramOrder {
		if _, ok := inputs[p]; ok {
			keys = append(keys, p)
			used[p] = true
		}
	}
	for k := range inputs {
		if !used[k] {
			keys = append(keys, k)
		}
	}
	return keys
}
